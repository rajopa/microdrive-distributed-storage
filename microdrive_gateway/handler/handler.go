package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	pb "microdrive_gateway/pkg/proto"
)

type Handler struct {
	authClient    pb.AuthClient
	storageClient pb.ImageServiceClient
	paymentClient pb.PaymentServiceClient
}

func NewHandler(auth pb.AuthClient, storage pb.ImageServiceClient, payment pb.PaymentServiceClient) *Handler {
	return &Handler{
		authClient:    auth,
		storageClient: storage,
		paymentClient: payment,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/sign-in", h.SignIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		api.POST("/upload", h.UploadImage)
		api.POST("/pay", h.ProcessPayment)
	}
	return router
}

func (h *Handler) SignUp(c *gin.Context) {
    var req pb.RegisterRequest
    if err := c.BindJSON(&req); err != nil {
        c.AbortWithStatusJSON(400, gin.H{"error": "invalid body"})
        return
    }

    res, err := h.authClient.Register(c.Request.Context(), &req)
    if err != nil {
        c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{
        "user_id": res.UserId,
        "email":   req.Email,      
        "message": "user created successfully",
        "status":  "active",
    })
}
func (h *Handler) SignIn(c *gin.Context) {
	var req pb.LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid body"})
		return
	}

	// Добавил AppId, так как он есть в твоем proto
	req.AppId = 1

	res, err := h.authClient.Login(c.Request.Context(), &req)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.AbortWithStatusJSON(401, gin.H{"error": "empty auth header"})
		return
	}

	// В твоем auth.proto нет метода проверки токена.
	// Пока просто пропускаем, чтобы не было ошибок.
	// res, err := h.authClient.IsAdmin(c.Request.Context(), &pb.IsAdminRequest{UserId: ...})

	c.Set("userId", int64(1))
}

func (h *Handler) UploadImage(c *gin.Context) {

}

func (h *Handler) ProcessPayment(c *gin.Context) {
	var req pb.PaymentRequest
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid body"})
		return
	}

	res, err := h.paymentClient.ProcessPayment(c.Request.Context(), &req)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}
