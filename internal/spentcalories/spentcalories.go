package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("incorrect input")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("the numbers of steps must be greater than zero")
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("the duration of the walk should be positive")
	}

	return steps, parts[1], duration, nil
}

func distance(steps int, height float64) float64 {
	passed := ((height * stepLengthCoefficient) * float64(steps)) / mInKm
	return passed
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	passedDistance := distance(steps, height)

	durationHours := duration.Hours()
	if durationHours == 0 {
		return 0
	}

	return passedDistance / durationHours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	durationHours := duration.Hours()
	var calories float64

	switch strings.ToLower(activity) {
	case "бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		strings.Title(activity), durationHours, dist, speed, calories,
	)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("incorrect input")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	calories := (durationMinutes * speed * weight) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("all entered data must be greater than zero")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	calories := ((durationMinutes * speed * weight) / minInH) * walkingCaloriesCoefficient

	return calories, nil
}
