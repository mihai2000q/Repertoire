import { ActionIcon, Flex, Grid, Stack, Text, Tooltip } from '@mantine/core'
import { IconCheck } from '@tabler/icons-react'
import dayjs from 'dayjs'
import EditPanelCard from '../../@ui/card/EditPanelCard.tsx'
import Song from '../../../types/models/Song.ts'
import { useDisclosure } from '@mantine/hooks'
import EditSongInformationModal from '../modal/EditSongInformationModal.tsx'
import DifficultyBar from '../../@ui/bar/DifficultyBar.tsx'

const NotSet = ({ label }: { label?: string }) => (
  <Text fz={'sm'} c={'dimmed'} fs={'oblique'} inline>
    {label ? label : 'not set'}
  </Text>
)

const firstColSize = 6
const secondColSize = 6

interface SongInformationCardProps {
  song: Song
}

function SongInformationCard({ song }: SongInformationCardProps) {
  const [openedEdit, { open: openEdit, close: closeEdit }] = useDisclosure(false)

  return (
    <EditPanelCard p={'md'} ariaLabel={'song-information-card'} onEditClick={openEdit}>
      <Stack gap={'xs'}>
        <Text fw={600}>Information</Text>
        <Grid align={'center'} gutter={'sm'}>
          <Grid.Col span={firstColSize}>
            <Text fw={500} c={'dimmed'}>
              Difficulty:
            </Text>
          </Grid.Col>
          <Grid.Col span={secondColSize}>
            {song.difficulty ? (
              <DifficultyBar difficulty={song.difficulty} maw={160} />
            ) : (
              <NotSet />
            )}
          </Grid.Col>

          <Grid.Col span={firstColSize}>
            <Text fw={500} c={'dimmed'} truncate={'end'}>
              Guitar Tuning:
            </Text>
          </Grid.Col>
          <Grid.Col span={secondColSize}>
            {song.guitarTuning ? <Text fw={600}>{song.guitarTuning.name}</Text> : <NotSet />}
          </Grid.Col>

          <Grid.Col span={firstColSize}>
            <Flex>
              <Tooltip label={'Beats Per Minute'} openDelay={200} position={'top-start'}>
                <Text fw={500} c={'dimmed'}>
                  Bpm:
                </Text>
              </Tooltip>
            </Flex>
          </Grid.Col>
          <Grid.Col span={secondColSize}>
            {song.bpm ? <Text fw={600}>{song.bpm}</Text> : <NotSet />}
          </Grid.Col>

          <Grid.Col span={firstColSize}>
            <Text fw={500} c={'dimmed'}>
              Recorded:
            </Text>
          </Grid.Col>
          <Grid.Col span={secondColSize}>
            {song.isRecorded ? (
              <ActionIcon
                component={'div'}
                size={'20px'}
                aria-label={'recorded-icon'}
                sx={(theme) => ({
                  cursor: 'default',
                  backgroundColor: theme.colors.primary[5],
                  '&:hover': { backgroundColor: theme.colors.primary[5] },
                  '&:active': { transform: 'none' }
                })}
              >
                <IconCheck size={14} />
              </ActionIcon>
            ) : (
              <Text fw={600}>No</Text>
            )}
          </Grid.Col>

          <Grid.Col span={firstColSize}>
            <Flex>
              <Tooltip
                label={"This field is calculated based on sections' rehearsals"}
                openDelay={200}
                position={'top-start'}
              >
                <Text fw={500} c={'dimmed'}>
                  Last Played On:
                </Text>
              </Tooltip>
            </Flex>
          </Grid.Col>

          <Grid.Col span={secondColSize}>
            {song.lastTimePlayed ? (
              <Flex>
                <Tooltip
                  label={`Song was played last time on ${dayjs(song.lastTimePlayed).format('D MMMM YYYY [at] hh:mm A')}`}
                  openDelay={400}
                  disabled={!song.lastTimePlayed}
                >
                  <Text fw={600}>{dayjs(song.lastTimePlayed).format('DD MMM YYYY')}</Text>
                </Tooltip>
              </Flex>
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
