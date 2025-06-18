package daysteps

import (
	"fmt"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	var steps int
	var duration time.Duration
	var err error
	entries := strings.Split(data, ",")
	if len(entries) != 2 {
		return 0, 0, fmt.Errorf("неверный формат")
	}
	steps, err = strconv.Atoi(entries[0])
	if err != nil {
		return 0, 0, err
	}
	if steps < 1 {
		return 0, 0, fmt.Errorf("неверное число шагов")
	}
	duration, err = time.ParseDuration(entries[1])
	if err != nil {
		return 0, 0, err
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("неверная продолжительность")
	}
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	if steps < 1 {
		return ""
	}
	distance := stepLength * float64(steps) / mInKm
	kcal, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}
	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps,
		distance,
		kcal,
	)
}
