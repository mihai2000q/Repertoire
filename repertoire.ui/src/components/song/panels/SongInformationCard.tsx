import { ActionIcon, Grid, Group, Progress, Stack, Text, Tooltip } from '@mantine/core'
import Difficulty from '../../../utils/enums/Difficulty.ts'
import { IconCheck } from '@tabler/icons-react'
import dayjs from 'dayjs'
import EditPanelCard from '../../card/EditPanelCard.tsx'
import Song from '../../../types/models/Song.ts'
import { useDisclosure } from '@mantine/hooks'
import useDifficultyInfo from '../../../hooks/useDifficultyInfo.ts'
import EditSongInformationModal from '../modal/EditSongInformationModal.tsx'

const NotSet = ({ label }: { label?: string }) => (
  <Text fz={'sm'} c={'dimmed'} fs={'oblique'} inline>
    {label ? label : 'not set'}
  </Text>
)

interface SongInformationCardProps {
  song: Song
}

function SongInformationCard({ song }: SongInformationCardProps) {
  const [openedEdit, { open: openEdit, close: closeEdit }] = useDisclosure(false)

  const { number: difficultyNumber, color: difficultyColor } = useDifficultyInfo(song?.difficulty)

  return (
    <EditPanelCard p={'md'} onEditClick={openEdit}>
      <Stack gap={'xs'}>
        <Text fw={600}>Information</Text>
        <Grid align={'center'} gutter={'sm'}>
          <Grid.Col span={6}>
            <Text fw={500} c={'dimmed'}>
              Difficulty:
            </Text>
          </Grid.Col>
          <Grid.Col span={6}>
            {song.difficulty ? (
              <Tooltip label={`This song is ${song.difficulty}`} openDelay={400} position={'top'}>
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
            ) : (
              <NotSet />
            )}
          </Grid.Col>

          <Grid.Col span={6}>
            <Text fw={500} c={'dimmed'} truncate={'end'}>
              Guitar Tuning:
            </Text>
          </Grid.Col>
          <Grid.Col span={6}>
            {song.guitarTuning ? <Text fw={600}>{song.guitarTuning.name}</Text> : <NotSet />}
          </Grid.Col>

          <Grid.Col span={6}>
            <Tooltip label={'Beats Per Minute'} openDelay={200} position={'top-start'}>
              <Text fw={500} c={'dimmed'}>
                Bpm:
              </Text>
            </Tooltip>
          </Grid.Col>
          <Grid.Col span={6}>{song.bpm ? <Text fw={600}>{song.bpm}</Text> : <NotSet />}</Grid.Col>

          <Grid.Col span={6}>
            <Text fw={500} c={'dimmed'}>
              Recorded:
            </Text>
          </Grid.Col>
          <Grid.Col span={6}>
            {song.isRecorded ? (
              <ActionIcon
                size={'20px'}
                sx={(theme) => ({
                  cursor: 'default',
                  backgroundColor: theme.colors.cyan[5],
                  '&:hover': { backgroundColor: theme.colors.cyan[5] }
                })}
              >
                <IconCheck size={14} />
              </ActionIcon>
            ) : (
              <Text fw={600}>No</Text>
            )}
          </Grid.Col>

          <Grid.Col span={6}>
            <Tooltip
              label={"This field is calculated based on sections' rehearsals"}
              openDelay={200}
              position={'top-start'}
            >
              <Text fw={500} c={'dimmed'}>
                Last Played On:
              </Text>
            </Tooltip>
          </Grid.Col>

          <Grid.Col span={6}>
            {song.lastTimePlayed ? (
              <Text fw={600}>{dayjs(song.lastTimePlayed).format('DD MMM YYYY')}</Text>
            ) : (
              <NotSet label={'never'} />
            )}
          </Grid.Col>
        </Grid>
      </Stack>

      <EditSongInformationModal song={song} opened={openedEdit} onClose={closeEdit} />
    </EditPanelCard>
  )
}

export default SongInformationCard
