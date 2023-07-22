package helper

import "errors"

// error message conventions
var (
	// ErrAuthenticationFailed error wrong authentication data
	ErrAuthenticationFailed = errors.New("Email atau kata sandi salah")

	// ErrUserNotFound error user does not exist
	ErrUserNotFound = errors.New("Pengguna tidak ditemukan")

	// ErrMenteeNotFound error mentee does not exist
	ErrMenteeNotFound = errors.New("Mentee tidak ditemukan")

	// ErrMentorNotFound error mentor does not exist
	ErrMentorNotFound = errors.New("Mentor tidak ditemukan")

	// ErrCategoryNotFound error category does not exist
	ErrCategoryNotFound = errors.New("Kategori tidak ditemukan")

	// ErrCourseNotFound error course does not exist
	ErrCourseNotFound = errors.New("Kursus tidak ditemukan")

	// ErrModuleNotFound error module does not exist
	ErrModuleNotFound = errors.New("Modul tidak ditemukan")

	// ErrMaterialNotFound error material does not exist
	ErrMaterialNotFound = errors.New("Materi tidak ditemukan")

	// ErrMaterialAssetNotFound error material asset does not exist
	ErrMaterialAssetNotFound = errors.New("Aset materi tidak ditemukan")

	// ErrAssignmentNotFound error assignment does not exist
	ErrAssignmentNotFound = errors.New("Tugas tidak ditemukan")

	// ErrAssignmentNotFound error assignment does not exist
	ErrAssignmentMenteeNotFound = errors.New("Tugas mentee tidak ditemukan")

	// ErrEmailAlreadyExist error email already exist
	ErrEmailAlreadyExist = errors.New("Email telah digunakan")

	// ErrPasswordLengthInvalid error invalid password length
	ErrPasswordLengthInvalid = errors.New("Panjang password minimal 6 karakter")

	// ErrPasswordNotMatch error both password not match
	ErrPasswordNotMatch = errors.New("Kedua password tidak sama")

	// ErrOTPExpired error otp expired
	ErrOTPExpired = errors.New("OTP telah kadaluarsa")

	// ErrOTPNotMatch error OTP not match with server
	ErrOTPNotMatch = errors.New("OTP yang anda masukkan salah")

	// ErrAccessForbidden error access forbidden
	ErrAccessForbidden = errors.New("Akses dilarang")

	// ErrUserUnauthorized error user unauthorized
	ErrUserUnauthorized = errors.New("User tidak ter-Autentikasi")

	// ErrInvalidRequest error invalid request body
	ErrInvalidRequest = errors.New("Invalid request body")

	// ErrInvalidJWTPayload error invalid JWT payloads
	ErrInvalidJWTPayload = errors.New("Invalid JWT payloads")

	// ErrUnsupportedAssignmentFile error unsupported file upload
	ErrUnsupportedAssignmentFile = errors.New("Extensi file tugas tidak didukung. Gunakan file ber-extensi .pdf")

	// ErrInvalidTokenHeader error invalid token header
	ErrInvalidTokenHeader = errors.New("Invalid token header")

	// ErrUnsupportedVideoFile error unsupported video file
	ErrUnsupportedVideoFile = errors.New("Extensi file video tidak didukung. Gunakan file ber-extensi .mp4 atau .mkv")

	// ErrUnsupportedImageFile error unsupported image file
	ErrUnsupportedImageFile = errors.New("Extensi file gambar tidak didukung. Gunakan file ber-extensi .jpeg, .jpg, atau .png")

	// ErrRecordNotFound error record not found (cannot specify the error)
	ErrRecordNotFound = errors.New("Data tidak ditemukan")

	// ErrAlreadyEnrolled error already enrolled this course
	ErrAlreadyEnrolled = errors.New("Anda telah mengambil kursus ini")

	// ErrNoEnrolled error no enrolled this course
	ErrNoEnrolled = errors.New("Kamu harus mengambil kursus ini terlebih dahulu")

	// ErrInternalServerError error internal server error
	ErrInternalServerError = errors.New("Internal server error")
)
