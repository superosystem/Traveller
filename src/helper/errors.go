package helper

import "errors"

// error message conventions
var (
	// ErrAuthenticationFailed error wrong authentication data
	ErrAuthenticationFailed = errors.New("email atau kata sandi salah")

	// ErrUserNotFound error user does not exist
	ErrUserNotFound = errors.New("pengguna tidak ditemukan")

	// ErrMenteeNotFound error mentee does not exist
	ErrMenteeNotFound = errors.New("mentee tidak ditemukan")

	// ErrMentorNotFound error mentor does not exist
	ErrMentorNotFound = errors.New("mentor tidak ditemukan")

	// ErrCategoryNotFound error category does not exist
	ErrCategoryNotFound = errors.New("kategori tidak ditemukan")

	// ErrCourseNotFound error course does not exist
	ErrCourseNotFound = errors.New("kursus tidak ditemukan")

	// ErrModuleNotFound error module does not exist
	ErrModuleNotFound = errors.New("modul tidak ditemukan")

	// ErrMaterialNotFound error material does not exist
	ErrMaterialNotFound = errors.New("materi tidak ditemukan")

	// ErrMaterialAssetNotFound error material asset does not exist
	ErrMaterialAssetNotFound = errors.New("aset materi tidak ditemukan")

	// ErrAssignmentNotFound error assignment does not exist
	ErrAssignmentNotFound = errors.New("tugas tidak ditemukan")

	// ErrAssignmentNotFound error assignment does not exist
	ErrAssignmentMenteeNotFound = errors.New("tugas mentee tidak ditemukan")

	// ErrEmailAlreadyExist error email already exist
	ErrEmailAlreadyExist = errors.New("email telah digunakan")

	// ErrPasswordLengthInvalid error invalid password length
	ErrPasswordLengthInvalid = errors.New("panjang password minimal 6 karakter")

	// ErrPasswordNotMatch error both password not match
	ErrPasswordNotMatch = errors.New("kedua password tidak sama")

	// ErrOTPExpired error otp expired
	ErrOTPExpired = errors.New("otp telah kadaluarsa")

	// ErrOTPNotMatch error OTP not match with server
	ErrOTPNotMatch = errors.New("otp yang anda masukkan salah")

	// ErrAccessForbidden error access forbidden
	ErrAccessForbidden = errors.New("akses dilarang")

	// ErrUserUnauthorized error user unauthorized
	ErrUserUnauthorized = errors.New("user tidak ter-Autentikasi")

	// ErrInvalidRequest error invalid request body
	ErrInvalidRequest = errors.New("invalid request body")

	// ErrInvalidJWTPayload error invalid JWT payloads
	ErrInvalidJWTPayload = errors.New("invalid JWT payloads")

	// ErrUnsupportedAssignmentFile error unsupported file upload
	ErrUnsupportedAssignmentFile = errors.New("extensi file tugas tidak didukung. Gunakan file ber-extensi .pdf")

	// ErrInvalidTokenHeader error invalid token header
	ErrInvalidTokenHeader = errors.New("invalid token header")

	// ErrUnsupportedVideoFile error unsupported video file
	ErrUnsupportedVideoFile = errors.New("extensi file video tidak didukung. Gunakan file ber-extensi .mp4 atau .mkv")

	// ErrUnsupportedImageFile error unsupported image file
	ErrUnsupportedImageFile = errors.New("extensi file gambar tidak didukung. Gunakan file ber-extensi .jpeg, .jpg, atau .png")

	// ErrRecordNotFound error record not found (cannot specify the error)
	ErrRecordNotFound = errors.New("data tidak ditemukan")

	// ErrAlreadyEnrolled error already enrolled this course
	ErrAlreadyEnrolled = errors.New("anda telah mengambil kursus ini")

	// ErrNoEnrolled error no enrolled this course
	ErrNoEnrolled = errors.New("kamu harus mengambil kursus ini terlebih dahulu")

	// ErrInternalServerError error internal server error
	ErrInternalServerError = errors.New("internal server error")
)
