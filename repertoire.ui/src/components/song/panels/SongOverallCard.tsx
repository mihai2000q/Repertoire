import { Card, Flex, Grid, NumberFormatter, Stack, Text, Tooltip } from '@mantine/core'
import Song from '../../../types/models/Song.ts'
import SongProgressBar from '../../@ui/misc/SongProgressBar.tsx'
import SongConfidenceBar from '../../@ui/misc/SongConfidenceBar.tsx'

interface SongOverallCardProps {
  song: Song
}

function SongOverallCard({ song }: SongOverallCardProps) {
  return (
    <Card variant={'panel'} aria-label={'song-overall-card'} p={'md'}>
      <Stack>
        <Tooltip
          label={"This panel is calculated based on sections' data"}
          openDelay={300}
          position={'top-start'}
        >
          <Text fw={600} style={{ alignSelf: 'start' }}>
            Overall
          </Text>
        </Tooltip>

        <Grid align={'center'} gutter={'sm'}>
          <Grid.Col span={6}>
            <Flex>
              <Tooltip
                label={'Median of total number of rehearsals'}
                openDelay={200}
                position={'top-start'}
              >
                <Text fw={500} c={'dimmed'} truncate={'end'}>
                  Rehearsals:
                </Text>
              </Tooltip>
            </Flex>
          </Grid.Col>
          <Grid.Col span={6}>
            <Text fw={600}>
              <NumberFormatter value={song.rehearsals} />
            </Text>
          </Grid.Col>

          <Grid.Col span={6}>
            <Flex>
              <Tooltip label={'Median of confidence'} openDelay={200} position={'top-start'}>
                <Text fw={500} c={'dimmed'} truncate={'end'}>
                  Confidence:
                </Text>
              </Tooltip>
            </Flex>
          </Grid.Col>
          <Grid.Col span={6}>
            <SongConfidenceBar confidence={song.confidence} flex={1} />
          </Grid.Col>

          <Grid.Col span={6}>
            <Flex>
              <Tooltip label={'Median of progress'} openDelay={200} position={'top-start'}>
                <Text fw={500} c={'dimmed'} truncate={'end'}>
                  Progress:
                </Text>
              </Tooltip>
            </Flex>
          </Grid.Col>
          <Grid.Col span={6}>
            <SongProgressBar progress={song.progress} flex={1} />
          </Grid.Col>
        </Grid>
      </Stack>
    </Card>
  )
}

export default SongOverallCard
