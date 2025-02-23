package study

import (
	"errors"
	"fmt"
	"github.com/szwtdl/jxtxedu/types"
	"github.com/szwtdl/jxtxedu/utils"
	"github.com/szwtdl/req"
)

// 获取课程列表

func CourseList(client *client.HttpClient) ([]types.Course, error) {
	// 请求数据并解析响应
	response, err := client.DoGet(types.CourseListUrl)
	if err != nil {
		return nil, err
	}

	var responseApi types.ResponseApi
	if err := utils.JsonUnmarshal(response, &responseApi); err != nil {
		return nil, err
	}

	// 检查返回代码
	if responseApi.Code != 0 {
		return nil, errors.New(responseApi.Msg)
	}

	// 获取并解密数据
	decryptData, err := utils.DecryptEncryptedData(responseApi)
	if err != nil {
		return nil, err
	}
	// 解析解密后的课程数据并返回分析后的课程列表
	return analyzeData(decryptData)
}

// 获取课程章节ID

func GetLessonId(client *client.HttpClient, courseId int) (int, error) {
	// 请求数据并解析响应
	response, err := client.DoGet(fmt.Sprintf("%s?course_id=%d", types.ChapterUrl, courseId))
	if err != nil {
		return 0, err
	}
	var responseApi types.ResponseApi
	if err := utils.JsonUnmarshal(response, &responseApi); err != nil {
		return 0, err
	}
	// 检查返回代码
	if responseApi.Code != 0 {
		return 0, errors.New(responseApi.Msg)
	}
	decryptData, err := utils.DecryptEncryptedData(responseApi)
	if err != nil {
		return 0, err
	}
	// 获取并解密数据
	var Items struct {
		IsBuy            int `json:"is_buy"`
		LearningLessonId int `json:"learning_lesson_id"`
	}
	if err := utils.JsonUnmarshal([]byte(decryptData), &Items); err != nil {
		return 0, err
	}
	if Items.IsBuy == 0 {
		return 0, errors.New("未购买课程")
	}
	return Items.LearningLessonId, nil
}

// 课程学习位置

func ChapterProgress(client *client.HttpClient, courseId int, lessonId int) (types.ChapterProgress, error) {
	response, err := client.DoGet(fmt.Sprintf("%s?course_id=%d&lesson_id=%d", types.VideoAuthUrl, courseId, lessonId))
	if err != nil {
		return types.ChapterProgress{}, err
	}
	var responseApi types.ResponseApi
	if err := utils.JsonUnmarshal(response, &responseApi); err != nil {
		return types.ChapterProgress{}, err
	}
	if responseApi.Code != 0 {
		return types.ChapterProgress{}, errors.New(responseApi.Msg)
	}

	decryptData, err := utils.DecryptEncryptedData(responseApi)
	if err != nil {
		return types.ChapterProgress{}, err
	}
	var courseInfo types.CourseInfo
	if err := utils.JsonUnmarshal([]byte(decryptData), &courseInfo); err != nil {
		return types.ChapterProgress{}, err
	}
	var chapterProgress types.ChapterProgress
	chapterProgress.Title = courseInfo.Video.Title
	finishSecond, _ := utils.ToInt(courseInfo.Record.FinishSecond)
	startSecond, _ := utils.ToInt(courseInfo.Record.StartSecond)
	chapterProgress.FinishSecond = finishSecond
	chapterProgress.StartSecond = startSecond
	chapterProgress.Face = make([]int, len(courseInfo.Face))
	for i, faceStr := range courseInfo.Face {
		floatValue, _ := utils.ToInt(faceStr)
		chapterProgress.Face[i] = floatValue
	}
	chapterProgress.IsDone = courseInfo.Lesson.IsDone
	chapterProgress.Duration = courseInfo.Video.Duration
	chapterProgress.Photo = courseInfo.User.Photo
	return chapterProgress, nil
}

// 分析并返回课程数据
func analyzeData(decryptData string) ([]types.Course, error) {
	var courseResponse types.CourseResponse
	if err := utils.JsonUnmarshal([]byte(decryptData), &courseResponse); err != nil {
		return nil, err
	}
	// 使用预分配容量，减少内存扩容操作
	courses := make([]types.Course, len(courseResponse.Data))
	for i, item := range courseResponse.Data {
		courses[i] = types.Course{
			Id:                    item.Id,
			Title:                 item.Title,
			Cover:                 item.Cover,
			CourseGraduationTotal: item.CourseGraduationTotal,
			TotalLearningTime:     item.TotalLearningTime,
			LearningProgress:      item.LearningProgress,
			DurationCount:         item.DurationCount,
			TeacherName:           item.TeacherName,
			TeacherAvatar:         item.TeacherAvatar,
			Progress:              item.Progress,
			OrderStatus:           item.OrderStatus,
			ChapterCount:          item.ChapterCount,
			LessonCount:           item.LessonCount,
			CompletedHour:         item.CompletedHour,
		}
	}
	return courses, nil
}
