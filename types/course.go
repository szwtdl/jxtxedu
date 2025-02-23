package types

type Course struct {
	Id                    int         `json:"id"`
	Title                 string      `json:"title"`
	Cover                 string      `json:"cover"`
	TeacherName           string      `json:"teacher_name"`
	TeacherAvatar         string      `json:"teacher_avatar"`
	OrderStatus           int         `json:"order_status"`
	CourseGraduationTotal int         `json:"course_graduation_total"`
	TotalLearningTime     int         `json:"total_learning_time"`
	LearningProgress      float64     `json:"learning_progress"`
	Progress              interface{} `json:"progress"`
	DurationCount         int         `json:"duration_count"`
	CompletedHour         interface{} `json:"completed_hour"`
	ChapterCount          int         `json:"chapter_count"`
	LessonCount           int         `json:"lesson_count"`
	Status                bool        `json:"status"`
}

type CourseInfo struct {
	Video struct {
		Id       string `json:"id"`
		Title    string `json:"title"`
		Cover    string `json:"cover"`
		Duration int    `json:"duration"`
		Player   struct {
			All    string `json:"all"`
			Tx     string `json:"tx"`
			Baidu  string `json:"baidu"`
			Huawei string `json:"huawei"`
			Dx     string `json:"dx"`
		} `json:"player"`
		Source int `json:"source"`
	} `json:"video"`
	Lesson struct {
		Id         int    `json:"id"`
		Title      string `json:"title"`
		IsFree     bool   `json:"is_free"`
		IsDone     bool   `json:"is_done"`
		IsForward  bool   `json:"is_forward"`
		IsInExam   int    `json:"is_in_exam"`
		ExamTimes  int    `json:"exam_times"`
		FreeSecond int    `json:"free_second"`
		IsPractice bool   `json:"is_practice"`
	}
	Record struct {
		StartSecond  int         `json:"start_second"`
		FinishSecond interface{} `json:"finish_second"`
		PassScore    int         `json:"pass_score"`
	} `json:"record"`
	Face       []interface{} `json:"face"`
	FaceConfig struct {
		Type int `json:"type"`
	}
	User struct {
		RealStatus bool `json:"real_status"`
		Photo      int  `json:"photo"`
	}
}

type ChapterProgress struct {
	Title        string `json:"title"`
	FinishSecond int    `json:"finish_second"`
	StartSecond  int    `json:"start_second"`
	Duration     int    `json:"duration"`
	IsDone       bool   `json:"is_done"`
	Face         []int  `json:"face"`
	Photo        int    `json:"photo"`
}
