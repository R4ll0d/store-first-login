# Stage 1: Build Stage
FROM golang:1.20-alpine AS builder

# กำหนด working directory
WORKDIR /app/store-first-login

# คัดลอก go.mod และ go.sum
COPY go.mod go.sum ./

# ติดตั้ง dependencies
RUN go mod tidy

# คัดลอกโค้ดทั้งหมด
COPY . .

# สร้างโปรแกรม
RUN go build -ldflags "-X main.Buildtime=$(date +%FT%T%z)" -o store-first-login main.go

# Stage 2: Run Stage
FROM alpine:latest

# ติดตั้ง dependencies 
RUN go mod tidy

# กำหนด working directory
WORKDIR /app/store-first-login

# คัดลอกไฟล์ที่ build จาก stage แรก
COPY --from=builder /app/store-first-login .

# รันแอปพลิเคชัน
CMD ["./myapp"]
