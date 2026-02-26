package vendor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	vendorService Service
}

func NewHandler(vendorService Service) *Handler {
	return &Handler{
		vendorService: vendorService,
	}
}

func (h *Handler) CompleteVendorProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	var req CompleteVendorRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	vendorProfile, err := h.vendorService.CompleteVendorProfile(c.Request.Context(), userID.(primitive.ObjectID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,		 gin.H{"error": err.Error()})
		return
	}			
	c.JSON(http.StatusOK, vendorProfile)
}				

func (h *Handler) GetVendorProfile(c *gin.Context) {
	userID, exists := c.Get("userID")			
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	vendorProfile, err := h.vendorService.GetVendorProfile(c.Request.Context(), userID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vendorProfile)
}

func (h *Handler) UpdateVendorProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	var req UpdateVendorProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	vendorProfile, err := h.vendorService.UpdateVendorProfile(c.Request.Context(), userID.(primitive.ObjectID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vendorProfile)
}

func (h *Handler) DeactivateVendorProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	err := h.vendorService.DeactivateVendorProfile(c.Request.Context(), userID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vendor profile deactivated"})
}	
