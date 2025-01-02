import useDifficultyInfo from './useDifficultyInfo.ts'
import Difficulty from '../utils/enums/Difficulty.ts'
import { MantineTheme, useMantineTheme } from '@mantine/core'
import { mantineRenderHook } from '../test-utils.tsx'

describe('use Difficulty Info', () => {
  it.each([
    [Difficulty.Easy, { number: 1, getColor: (theme: MantineTheme) => theme.colors.green[5] }],
    [Difficulty.Medium, { number: 2, getColor: (theme: MantineTheme) => theme.colors.yellow[5] }],
    [Difficulty.Hard, { number: 3, getColor: (theme: MantineTheme) => theme.colors.orange[5] }],
    [Difficulty.Impossible, { number: 4, getColor: (theme: MantineTheme) => theme.colors.red[6] }]
  ])('should return difficulty info based on difficulty', (difficulty, output) => {
    // Arrange - parameterized
    // Act
    const { result } = mantineRenderHook(() => useDifficultyInfo(difficulty))

    const { result: theme } = mantineRenderHook(() => useMantineTheme())

    // Assert
    expect(result.current).toStrictEqual({
      number: output.number,
      color: output.getColor(theme.current as MantineTheme)
    })
  })

  it('should return default values when the difficulty is undefined', () => {
    const { result } = mantineRenderHook(() => useDifficultyInfo(undefined))

    expect(result.current).toStrictEqual({
      number: 0,
      color: ''
    })
  })
})
