import Difficulty from '../utils/enums/Difficulty.ts'
import { useMemo } from 'react'
import { useMantineTheme } from '@mantine/core'

export default function useDifficultyInfo(difficulty: Difficulty): {
  number: number
  color: string
} {
  const theme = useMantineTheme()

  return useMemo(() => {
    let number: number
    let color: string

    switch (difficulty) {
      case Difficulty.Easy:
        number = 1
        color = theme.colors.green[5]
        break
      case Difficulty.Medium:
        number = 2
        color = theme.colors.yellow[5]
        break
      case Difficulty.Hard:
        number = 3
        color = theme.colors.orange[5]
        break
      case Difficulty.Impossible:
        number = 4
        color = theme.colors.red[6]
        break
    }

    return { number, color }
  }, [difficulty])
}
