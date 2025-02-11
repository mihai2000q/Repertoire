import { Group, MantineStyleProps, Progress, Tooltip } from '@mantine/core'
import Difficulty from '../../../utils/enums/Difficulty.ts'
import useDifficultyInfo from '../../../hooks/useDifficultyInfo.ts'

interface DifficultyBarProps extends MantineStyleProps {
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
      <Group grow gap={'xxs'} role={'progressbar'} aria-label={'difficulty'}>
        {Array.from(Array(Object.entries(Difficulty).length)).map((_, i) => (
          <Progress
            key={i}
            size={size}
            value={i + 1 <= difficultyNumber ? 100 : 0}
            color={difficultyColor}
            {...props}
          />
        ))}
      </Group>
    </Tooltip>
  )
}

export default DifficultyBar
