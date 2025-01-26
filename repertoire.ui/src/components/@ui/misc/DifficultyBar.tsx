import { Group, Progress, Tooltip } from '@mantine/core'
import Difficulty from '../../../utils/enums/Difficulty.ts'
import useDifficultyInfo from '../../../hooks/useDifficultyInfo.ts'

interface DifficultyBarProps {
  difficulty: Difficulty | undefined
  size?: number | string
  maw?: number
}

function DifficultyBar({ difficulty, maw, size = 5 }: DifficultyBarProps) {
  const { number: difficultyNumber, color: difficultyColor } = useDifficultyInfo(difficulty)

  return (
    <Tooltip label={`This song is ${difficulty}`} openDelay={400} position={'top'}>
      <Group grow gap={4} role={'progressbar'} aria-label={'difficulty'}>
        {Array.from(Array(Object.entries(Difficulty).length)).map((_, i) => (
          <Progress
            key={i}
            size={size}
            maw={maw}
            value={i + 1 <= difficultyNumber ? 100 : 0}
            color={difficultyColor}
          />
        ))}
      </Group>
    </Tooltip>
  )
}

export default DifficultyBar
