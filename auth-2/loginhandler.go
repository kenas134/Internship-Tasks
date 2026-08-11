type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

func loginHandler(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := userStore.GetByEmail(req.Email)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Invalid credentials",
        })
        return
    }

    // Compare hashed password
    if err := bcrypt.CompareHashAndPassword(
        []byte(user.Password),
        []byte(req.Password),
    ); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Invalid credentials",
        })
        return
    }

    token, err := generateToken(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to generate token",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "type":  "Bearer",
        "expires_in": 86400, // 24 hours in seconds
    })
}