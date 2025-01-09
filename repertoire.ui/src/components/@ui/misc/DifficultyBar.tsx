import { Box, Group, Progress, Tooltip } from '@mantine/core'
import Difficulty from '../../../utils/enums/Difficulty.ts'
import useDifficultyInfo from '../../../hooks/useDifficultyInfo.ts'

interface DifficultyBarProps {
  difficulty: Difficulty | undefined
}

function DifficultyBar({ difficulty }: DifficultyBarProps) {
  const { number: difficultyNumber, color: difficultyColor } = useDifficultyInfo(difficulty)

  return (
    <Box role={'progressbar'} aria-label={'difficulty'}>
      <Tooltip label={`This song is ${difficulty}`} openDelay={400} position={'top'}>
        <Group grow gap={4}>
          {Array.from(Array(Object.entries(Difficulty).length)).map((_, i) => (
            <Progress
              key={i}
              size={5}
              maw={40}
              value={i + 1 <= difficultyNumber ? 100 : 0}
              color={difficultyColor}
            />
          ))}
        </Group>
      </Tooltip>
    </Box>
  )
}

export default DifficultyBar
