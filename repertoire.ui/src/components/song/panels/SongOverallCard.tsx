import { Card, Grid, NumberFormatter, Progress, Stack, Text, Tooltip } from '@mantine/core'
import Song from '../../../types/models/Song.ts'

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
          <Text fw={600}>Overall</Text>
        </Tooltip>

        <Grid align={'center'} gutter={'sm'}>
          <Grid.Col span={6}>
            <Tooltip
              label={'Median of total number of rehearsals'}
              openDelay={200}
              position={'top-start'}
            >
              <Text fw={500} c={'dimmed'} truncate={'end'}>
                Rehearsals:
              </Text>
            </Tooltip>
          </Grid.Col>
          <Grid.Col span={6}>
            <Text fw={600}>
              <NumberFormatter value={song.rehearsals} />
            </Text>
          </Grid.Col>

          <Grid.Col span={6}>
            <Tooltip label={'Median of confidence'} openDelay={200} position={'top-start'}>
              <Text fw={500} c={'dimmed'} truncate={'end'}>
                Confidence:
              </Text>
            </Tooltip>
          </Grid.Col>
          <Grid.Col span={6}>
            <Tooltip.Floating
              role={'tooltip'}
              label={
                <>
                  <NumberFormatter value={song.confidence} />%
                </>
              }
            >
              <Progress aria-label={'confidence'} flex={1} size={'sm'} value={song.confidence} />
            </Tooltip.Floating>
          </Grid.Col>

          <Grid.Col span={6}>
            <Tooltip label={'Median of progress'} openDelay={200} position={'top-start'}>
              <Text fw={500} c={'dimmed'} truncate={'end'}>
                Progress:
              </Text>
            </Tooltip>
          </Grid.Col>
          <Grid.Col span={6}>
            <Tooltip.Floating role={'tooltip'} label={<NumberFormatter value={song.progress} />}>
              <Progress
                aria-label={'progress'}
                flex={1}
                size={'sm'}
                value={song.progress / 10}
                color={'green'}
              />
            </Tooltip.Floating>
          </Grid.Col>
        </Grid>
      </Stack>
    </Card>
  )
}

export default SongOverallCard
