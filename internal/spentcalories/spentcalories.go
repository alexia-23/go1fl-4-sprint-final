package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	var steps int
	var activity string
	var duration time.Duration
	var err error
	entries := strings.Split(data, ",")
	if len(entries) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат")
	}
	steps, err = strconv.Atoi(entries[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps < 1 {
		return 0, "", 0, fmt.Errorf("неверные шаги")
	}
	activity = entries[1]
	duration, err = time.ParseDuration(entries[2])
	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("неверная длительность")
	}
	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	return float64(steps) * height * stepLengthCoefficient / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	dist := distance(steps, height)
	template := "Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n"
	var kcal float64
	switch activity {
	case "Бег":
		kcal, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			template,
			activity,
			duration.Hours(),
			dist,
			meanSpeed(steps, height, duration),
			kcal,
		), nil
	case "Ходьба":
		kcal, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			template,
			activity,
			duration.Hours(),
			dist,
			meanSpeed(steps, height, duration),
			kcal,
		), nil
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps < 1 {
		return 0, fmt.Errorf("неверное число шагов")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("неверный вес")
	}
	if height <= 0 {
		return 0, fmt.Errorf("неверная высота")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("неверная длительность")
	}
	speed := meanSpeed(steps, height, duration)
	return (weight * speed * duration.Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps < 1 {
		return 0, fmt.Errorf("неверное число шагов")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("неверный вес")
	}
	if height <= 0 {
		return 0, fmt.Errorf("неверная высота")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("неверная длительность")
	}
	speed := meanSpeed(steps, height, duration)
	return (weight * speed * duration.Minutes() * walkingCaloriesCoefficient) / minInH, nil
}
