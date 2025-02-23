package types

type ResponseApi struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

type EncryptedData struct {
	Ciphertext string `json:"ciphertext"`
	Salt       string `json:"salt"`
	IV         string `json:"iv"`
}

type CourseResponse struct {
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	Data        []struct {
		Id                     int         `json:"id"`
		Status                 int         `json:"status"`
		Title                  string      `json:"title"`
		Cover                  string      `json:"cover"`
		CertificateLessonCount int         `json:"certificate_lesson_count"`
		OnlineLessonCount      int         `json:"online_lesson_count"`
		TeacherName            string      `json:"teacher_name"`
		TeacherAvatar          string      `json:"teacher_avatar"`
		Expire                 int         `json:"expire"`
		DossierStatus          int         `json:"dossier_status"`
		DossierId              interface{} `json:"dossier_id"`
		DossierTime            interface{} `json:"dossier_time"`
		CompletedHour          interface{} `json:"completed_hour"`
		Progress               interface{} `json:"progress"`
		CourseGraduationTotal  int         `json:"course_graduation_total"`
		TotalLearningTime      int         `json:"total_learning_time"`
		LearningProgress       float64     `json:"learning_progress"`
		OrderStatus            int         `json:"order_status"`
		ChapterCount           int         `json:"chapter_count"`
		LessonCount            int         `json:"lesson_count"`
		DurationCount          int         `json:"duration_count"`
		ExamId                 int         `json:"exam_id"`
	} `json:"data"`
}
