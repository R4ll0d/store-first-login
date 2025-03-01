# Stage 1: Build Stage
FROM golang:1.23-alpine AS builder

# กำหนด working directory
WORKDIR /app

# คัดลอก go.mod และ go.sum
COPY go.mod go.sum ./

# ติดตั้ง dependencies
RUN go mod tidy

# คัดลอกโค้ดทั้งหมด
COPY . .

# สร้างโปรแกรม
RUN go build -o store-first-login main.go

# Stage 2: Run Stage
FROM alpine:latest

# กำหนด working directory
WORKDIR /app

# คัดลอกไฟล์ที่ build จาก stage แรก
COPY --from=builder /app/store-first-login .

# ระบุพอร์ตที่ต้องใช้
EXPOSE 8080

# รันแอปพลิเคชัน
CMD ["./store-first-login"]
