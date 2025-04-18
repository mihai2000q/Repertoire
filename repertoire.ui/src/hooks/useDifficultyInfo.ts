import Difficulty from '../types/enums/Difficulty.ts'
import { MantineTheme, useMantineTheme } from '@mantine/core'

interface useDifficultyInfoResult {
  number: number
  color: string
}

interface difficultyInfoMapResult {
  number: number
  getColor: (theme: MantineTheme) => string
}

const difficultyInfoMap = new Map<Difficulty, difficultyInfoMapResult>([
  [Difficulty.Easy, { number: 1, getColor: (theme) => theme.colors.green[5] }],
  [Difficulty.Medium, { number: 2, getColor: (theme) => theme.colors.yellow[5] }],
  [Difficulty.Hard, { number: 3, getColor: (theme) => theme.colors.orange[5] }],
  [Difficulty.Impossible, { number: 4, getColor: (theme) => theme.colors.red[6] }]
])

export default function useDifficultyInfo(difficulty: Difficulty): useDifficultyInfoResult {
  if (!difficulty) return { number: 0, color: '' }
  const theme = useMantineTheme()
  const result = difficultyInfoMap.get(difficulty)
  return {
    number: result.number,
    color: result.getColor(theme)
  }
}
