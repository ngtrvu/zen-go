package zen

var (
	// unauthorized
	ErrUnauthorized         = NewAppError("unauthorized", "unauthorized").AddTranslation("vi", "Không có quyền truy cập")
	ErrInvalidToken         = NewAppError("invalid_token", "invalid token").AddTranslation("vi", "Token không hợp lệ")
	ErrTokenExpired         = NewAppError("token_expired", "token expired").AddTranslation("vi", "Token hết hạn")
	ErrBadRequest           = NewAppError("400_bad_request", "bad request").AddTranslation("vi", "Yêu cầu không hợp lệ")
	ErrInvalidRequestFormat = NewAppError("invalid_request_format", "invalid request format").AddTranslation("vi", "Dữ liệu không hợp lệ")

	// forbidden
	ErrForbidden = NewAppError("forbidden", "forbidden").AddTranslation("vi", "Không có quyền truy cập")

	// not found
	ErrNotFound = NewAppError("not_found", "not found").AddTranslation("vi", "Không tìm thấy")

	// product
	ErrGetVcamProductFailed = NewAppError("get_vcam_product_failed", "get vcam product failed").AddTranslation("vi", "Lấy thông tin sản phẩm từ quỹ vcam thất bại. Vui lòng thử lại sau")
	ErrProductNotExists     = NewAppError("product_not_exists", "product not exists").AddTranslation("vi", "Sản phẩm không tồn tại")

	// incorrect input
	ErrInvalidImg       = NewAppError("invalid_image", "image is inapproriate").AddTranslation("vi", "Hình ảnh không hợp lệ")
	ErrIncorrectInput   = NewAppError("incorrect_input", "incorrect input or missing param").AddTranslation("vi", "Dữ liệu không hợp lệ")
	ErrInvalidImgFormat = NewAppError("invalid_image_format", "accept only image files in format JPEG/JPG/PNG").AddTranslation("vi", "Vui lòng sử dụng ảnh JPEG, JPG hoặc PNG")
	ErrOpenUploadFile   = NewAppError("open_upload_file_error", "open uploaded file error").AddTranslation("vi", "Lỗi khi mở file upload")
	ErrInvalidPhone     = NewAppError("invalid_phone", "incorrect phone format").AddTranslation("vi", "Số điện thoại không hợp lệ")
	ErrPhoneVerified    = NewAppError("phone_already_verified", "phone number is already verified").AddTranslation("vi", "Số điện thoại đã được xác thực")
	ErrIncorrectEmail   = NewAppError("incorrect_email", "email is missing or incorrect").AddTranslation("vi", "Email không hợp lệ")
	ErrEmailRegistered  = NewAppError("email_already_registered", "email is already registered").AddTranslation("vi", "Email đã được đăng ký")
	ErrAccountNotExist  = NewAppError("account_not_exist", "account does not exist in system").AddTranslation("vi", "Tài khoản không tồn tại")
	ErrInvalidLogin     = NewAppError("invalid_login", "email or password is incorrect").AddTranslation("vi", "Tên đăng nhập / Mật khẩu chưa chính xác, vui lòng thử lại")
	ErrIncorrectPwd     = NewAppError("incorrect_pwd", "incorrect password").AddTranslation("vi", "Mật khẩu không chính xác")
	ErrUnmatchedPwd     = NewAppError("unmatched_pwd", "confirm password does not match password").AddTranslation("vi", "Mật khẩu xác nhận không khớp")
	ErrIncorrectPIN     = NewAppError("incorrect_pin", "incorrect PIN").AddTranslation("vi", "PIN không chính xác")
	ErrUnmatchedPIN     = NewAppError("unmatched_pin", "confirm PIN does not match password").AddTranslation("vi", "PIN xác nhận không khớp")
	ErrFailedBiometrics = NewAppError("failed_biometrics", "failed FaceID/ Fingerprint ").AddTranslation("vi", "Xác thực khuôn mặt/ vân tay thất bại")
	ErrInvalidImgSize   = NewAppError("invalid_image_size", "ID card/ Selfie image file size exceeds maximum file size").AddTranslation("vi", "Kích thước ảnh vượt quá kích thước tối đa")
	ErrMinSellRequired  = NewAppError("min_sell_required", "sell Amount does not meet minimum sell").AddTranslation("vi", "Vui lòng bán tối thiểu 1 triệu đồng")

	// account
	ErrOperationFailedAccountApproveCanNotUpdateInfo = NewAppError("operation_failed_account_approved_can_not_update_info", "Operation failure because account is approved. Can not update info").AddTranslation("vi", "Tài khoản đã được duyệt. Không thể cập nhật thông tin")
	ErrOperationFailedFundAccountNotApprove          = NewAppError("operation_failed_fund_account_not_approve", "Operation failure because fund account is not approved").AddTranslation("vi", "Tài khoản quỹ chưa được duyệt")
	ErrOperationFailedFundOrderAlreadySubmit         = NewAppError("operation_failed_fund_order_already_submit", "Operation failure because fund order already submit").AddTranslation("vi", "Lệnh đã được gửi qua quỹ")
	ErrOperationFailedFundOrderSubmitFailed          = NewAppError("operation_failed_fund_order_submit_failed", "Operation failure because fund order submit failed").AddTranslation("vi", "Lệnh gửi qua quỹ thất bại. Vui lòng thử lại sau")
	ErrOperationFailedFundAccountExistsCanNotCreate  = NewAppError("operation_failed_fund_account_exists_can_not_create", "Operation failure because account exists can not create").AddTranslation("vi", "Tài khoản quỹ đã tồn tại. Không thể tạo mới. Vui lòng liên hệ Stag để được hỗ trợ.")
	ErrOperationFailedUserIdentificationNotFound     = NewAppError("operation_failed_user_identification_not_found", "Operation failure user identification not found").AddTranslation("vi", "Không tìm thấy thông tin eKYC. Vui lòng thực hiện định danh hoặc liên hệ Stag để được hỗ trợ.")
	ErrOperationFailedGetFrontIDImage                = NewAppError("operation_failed_get_front_id_image", "Operation failure get front ID image").AddTranslation("vi", "Không tìm thấy hình ảnh CMND mặt trước. Vui lòng thực hiện định danh hoặc liên hệ Stag để được hỗ trợ.")
	ErrOperationFailedGetBackIDImage                 = NewAppError("operation_failed_get_back_id_image", "Operation failure get back ID image").AddTranslation("vi", "Không tìm thấy hình ảnh CMND mặt sau. Vui lòng thực hiện định danh hoặc liên hệ Stag để được hỗ trợ.")
	ErrOperationFailedUserIdentificationInvalid      = NewAppError("operation_failed_user_identification_invalid", "Operation failure user identification invalid").AddTranslation("vi", "Thông tin định danh không hợp lệ. Vui lòng thực hiện định danh hoặc liên hệ Stag để được hỗ trợ.")
	ErrOperationFailedIDNumberExists                 = NewAppError("operation_failed_id_number_exists", "Operation failed because the ID number already exists").AddTranslation("vi", "Số CMND / CCCD đã tồn tại. Vui lòng kiểm tra nhập lại thông tin số CMND / CCCD hoặc liên hệ Stag theo địa chỉ support@stag.vn để được hỗ trợ.")

	// fund account
	ErrOperationFailedCreateVCAMAccount    = NewAppError("operation_failed_create_vcam_account", "Operation failure create vcam account").AddTranslation("vi", "Không thể tạo tài khoản công ty quản lý quỹ VCAM. Vui lòng thực hiện định danh hoặc liên hệ Stag để được hỗ trợ.")
	ErrOperationFailedOverwriteVCAMAccount = NewAppError("operation_failed_overwrite_vcam_account", "Operation failure overwrite vcam account").AddTranslation("vi", "Không thể ghi đè tài khoản công ty quản lý quỹ VCAM. Vui lòng thực hiện định danh hoặc liên hệ Stag để được hỗ trợ.")
	ErrOperationFailedGetVCAMAccount       = NewAppError("operation_failed_get_vcam_account", "Operation failure get vcam account").AddTranslation("vi", "Không thể lấy thông tin tài khoản công ty quản lý quỹ VCAM. Vui lòng thử lại.")
	ErrOperationFailedVCAMAccountNotExists = NewAppError("operation_failed_vcam_account_not_exists", "Operation failure vcam account not exists").AddTranslation("vi", "Không tìm thấy tài khoản công ty quản lý quỹ VCAM. Vui lòng thực hiện định danh hoặc liên hệ Stag để được hỗ trợ.")

	// order
	ErrInvalidQuantity                = NewAppError("invalid_quantity", "quantity must be greater than 0").AddTranslation("vi", "Số lượng phải lớn hơn 0")
	ErrMinBuyRequired                 = NewAppError("min_buy_required", "buy Amount does not meet min buy").AddTranslation("vi", "Vui lòng mua tối thiểu 1 triệu đồng")
	ErrOrderCompleted                 = NewAppError("order_completed", "order completed. can not process").AddTranslation("vi", "Lệnh đã hoàn tất. Không thể xử lý")
	ErrOrderCanceled                  = NewAppError("order_canceled", "order canceled. can not process").AddTranslation("vi", "Lệnh đã huỷ. Không thể xử lý")
	ErrGetVcamOrderFailed             = NewAppError("get_vcam_order_failed", "get vcam order failed").AddTranslation("vi", "Lấy thông tin lệnh từ quỹ vcam thất bại. Vui lòng thử lại sau")
	ErrGetVcamOrdersFailed            = NewAppError("get_vcam_orders_failed", "get vcam orders failed").AddTranslation("vi", "Lấy danh sách lệnh từ quỹ vcam thất bại. Vui lòng thử lại sau")
	ErrGetVcamOrderBuyEstimateFailed  = NewAppError("get_vcam_order_buy_estimate_failed", "get vcam order buy estimate failed").AddTranslation("vi", "Lấy thông tin ước tính mua từ quỹ vcam thất bại. Vui lòng thử lại sau")
	ErrGetVcamOrderSellEstimateFailed = NewAppError("get_vcam_order_sell_estimate_failed", "get vcam order sell estimate failed").AddTranslation("vi", "Lấy thông tin ước tính bán từ quỹ vcam thất bại. Vui lòng thử lại sau")
	ErrCancelVcamOrderFailed          = NewAppError("cancel_vcam_order_failed", "cancel vcam order failed").AddTranslation("vi", "Huỷ lệnh từ quỹ vcam thất bại. Vui lòng thử lại sau")
	ErrOTPExpired                     = NewAppError("expired_otp", "OTP is used or expired").AddTranslation("vi", "OTP đã được sử dụng hoặc hết hạn")
	ErrIncorrectOTP                   = NewAppError("incorrect_otp", "incorrect OTP").AddTranslation("vi", "OTP không chính xác")
	ErrOrderInvalid                   = NewAppError("order_invalid", "order invalid").AddTranslation("vi", "Lệnh không hợp lệ")
	ErrOrderNotSubmitToFund           = NewAppError("order_not_submit_to_fund", "order not submit to fund").AddTranslation("vi", "Lệnh chưa được gửi qua quỹ")
	ErrOrderSubmitFundFailed          = NewAppError("order_submit_to_fund_failed", "order submit to fund failed").AddTranslation("vi", "Lệnh gửi qua quỹ thất bại. Vui lòng thử lại sau")

	// portfolio input
	ErrSummarizePortfolio   = NewAppError("error_summarize_portfolio", "error summarize portfolio").AddTranslation("vi", "Không thể tổng hợp số dư")
	ErrInsufficientBalance  = NewAppError("insufficient_balance", "amount exceeds balance").AddTranslation("vi", "Số dư không đủ")
	ErrUnableGetBalance     = NewAppError("unable_get_balance", "unable gets balance").AddTranslation("vi", "Không thể kiểm tra số dư")
	ErrUnableGetBalanceFund = NewAppError("unable_get_balance_fund", "unable gets balance fund company").AddTranslation("vi", "Không thể kiểm tra số dư từ công ty quỹ")

	// user program
	ErrUserProgramCanNotDelete     = NewAppError("user_program_can_not_delete", "user program can not delete").AddTranslation("vi", "Chương trình đã có số dư, bạn không thể xoá. Vui lòng liên hệ Stag để được hỗ trợ.")
	ErrUserProgramNoActivePrograms = NewAppError("user_program_no_active_programs", "no user program active").AddTranslation("vi", "Không tìm thấy chương trình đang tham gia")
	ErrUserProgramNotFound         = NewAppError("user_program_not_found", "user program not found").AddTranslation("vi", "Không tìm thấy chương trình đang tham gia")

	// otp
	ErrQuotaExceeded = NewAppError("quota_exceeded", "Out of quota for repeated requests").AddTranslation("vi", "OTP đã hết lượt gửi. Vui lòng thử lại sau")

	// standard error
	ErrStandard = NewAppError("standard_error", "standard error").AddTranslation("vi", "Xử lý lỗi. Vui lòng thử lại sau")

	// employee
	ErrEmployeeNotOnboard = NewAppError("employee_not_onboard", "Employee not onboard").AddTranslation("vi", "Nhân viên chưa tham gia chương trình")
)
