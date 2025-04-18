import { Group, GroupProps, Progress, Tooltip } from '@mantine/core'
import Difficulty from '../../../types/enums/Difficulty.ts'
import useDifficultyInfo from '../../../hooks/useDifficultyInfo.ts'

interface DifficultyBarProps extends GroupProps {
  difficulty: Difficulty | undefined
  size?: number | string
}

function DifficultyBar({ difficulty, size = 5, ...props }: DifficultyBarProps) {
  const { number: difficultyNumber, color: difficultyColor } = useDifficultyInfo(difficulty)

  return (
    <Tooltip
      label={difficulty ? `This song is ${difficulty}` : 'This song has no difficulty set'}
      openDelay={400}
      position={'top'}
    >
      <Group
        grow
        gap={'xxs'}
        role={'progressbar'}
        aria-label={'difficulty'}
        wrap={'nowrap'}
        {...props}
      >
        {Array.from(Array(Object.entries(Difficulty).length)).map((_, i) => (
          <Progress
            key={i}
            size={size}
            value={i + 1 <= difficultyNumber ? 100 : 0}
            color={difficultyColor}
          />
        ))}
      </Group>
    </Tooltip>
  )
}

export default DifficultyBar
