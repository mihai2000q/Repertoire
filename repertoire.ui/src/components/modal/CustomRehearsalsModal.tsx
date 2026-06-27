import { useEffect, useState } from 'react'
import {
  alpha,
  Avatar,
  Button,
  Center,
  Group,
  Menu,
  Modal,
  ScrollArea,
  Skeleton,
  Stack,
  Text,
  Tooltip
} from '@mantine/core'
import useSearchBy from '../../hooks/api/useSearchBy.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'
import SongWithProperty from '../../types/enums/properties/SongWithProperty.ts'
import { useAddCustomSongRehearsalsMutation, useGetSongsQuery } from '../../state/api/songsApi.ts'
import Song, { SongArrangement } from '../../types/models/Song.ts'
import { AddCustomSongRehearsalRequest } from '../../types/requests/SongRequests.ts'
import { toast } from 'react-toastify'
import CustomIconMusicNoteEighth from '../icons/CustomIconMusicNoteEighth.tsx'
import { IconStarFilled } from '@tabler/icons-react'
import SongProperty from '../../types/enums/properties/SongProperty.ts'

interface CustomRehearsalsModalProps {
  opened: boolean
  onClose: () => void
  ids: string[]
  onSuccess?: () => void
}

function CustomRehearsalsModal({ opened, onClose, ids, onSuccess }: CustomRehearsalsModalProps) {
  const [addCustomRehearsals, { isLoading: isCustomLoading }] = useAddCustomSongRehearsalsMutation()

  const searchBy = useSearchBy([
    { property: SongProperty.Id, operator: FilterOperator.In, value: ids }
  ])
  const { data: songs, isLoading } = useGetSongsQuery(
    {
      searchBy: searchBy,
      with: [SongWithProperty.Arrangements]
    },
    { skip: !opened }
  )

  const [songsMap, setSongsMap] = useState<Map<string, Song>>(new Map<string, Song>())
  const [songArrangementsMap, setSongArrangementsMap] = useState<
    Map<string, SongArrangement | null>
  >(new Map<string, SongArrangement>())
  const [noArrangementsFound, setNoArrangementsFound] = useState(false)
  useEffect(() => {
    if (!songs) return
    const songsMap = new Map<string, Song>(songs.models.map((song) => [song.id, song]))
    setSongsMap(songsMap)

    const newSongArrangements = new Map<string, SongArrangement>()
    let arrangementsFound = 0
    ids.forEach((songId) => {
      const song = songsMap.get(songId)
      const arrangement =
        (song.arrangements.find((arrangement) => arrangement.id === song.defaultArrangementId) ??
        song.arrangements.length > 0)
          ? song.arrangements[0]
          : null
      newSongArrangements.set(songId, arrangement)
      if (arrangement !== null) arrangementsFound++
    })
    setSongArrangementsMap(newSongArrangements)
    setNoArrangementsFound(arrangementsFound === 0)
  }, [songs])

  function handleChangeArrangement(songId: string, arrangement: SongArrangement) {
    const newSongArrangements = new Map<string, SongArrangement>([...songArrangementsMap])
    newSongArrangements.set(songId, arrangement)
    setSongArrangementsMap(newSongArrangements)
  }

  async function handleAddCustomRehearsals() {
    const requests: AddCustomSongRehearsalRequest[] = Array.from(songArrangementsMap.entries())
      .filter(([_, arrangement]) => arrangement)
      .map(([songId, arrangement]) => ({
        id: songId,
        arrangementId: arrangement.id
      }))

    await addCustomRehearsals({ requests: requests }).unwrap()
    toast.success(`Custom rehearsals added!`)
    onClose()
    onSuccess?.()
  }

  if (isLoading || !songs)
    return (
      <Modal.Root opened={opened} onClose={onClose}>
        <Modal.Overlay />
        <Modal.Content>
          <Modal.Header>
            <Modal.Title>
              <Skeleton w={160} h={22} />
            </Modal.Title>
            <Modal.CloseButton />
          </Modal.Header>

          <Modal.Body>
            <Stack>
              <Stack px={'sm'} gap={'lg'}>
                {Array.from(Array(3)).map((_, i) => (
                  <Group key={i} justify={'space-between'}>
                    <Group gap={'xs'}>
                      <Skeleton w={38} h={38} />
                      <Stack gap={2}>
                        <Skeleton w={150} h={15} />
                        <Skeleton w={100} h={10} />
                      </Stack>
                    </Group>
                    <Skeleton w={100} h={20} />
                  </Group>
                ))}
              </Stack>

              <Skeleton w={'100%'} h={30} />
            </Stack>
          </Modal.Body>
        </Modal.Content>
      </Modal.Root>
    )

  return (
    <Modal opened={opened} onClose={onClose} title={'Custom Rehearsals'}>
      <Stack>
        <ScrollArea.Autosize mah={'50vh'} scrollbars={'y'} scrollbarSize={7}>
          <Stack gap={'lg'} px={'sm'}>
            {Array.from(songArrangementsMap.keys())
              .map((songId) => songsMap.get(songId))
              .map((song) => (
                <Group key={song.id} justify={'space-between'} wrap={'nowrap'}>
                  <Group gap={'xs'} wrap={'nowrap'}>
                    <Avatar
                      radius={'md'}
                      size={'md'}
                      src={song.imageUrl ?? song.album?.imageUrl}
                      alt={(song.imageUrl ?? song.album?.imageUrl) && song.title}
                      color={'gray.5'}
                    >
                      <Center c={'white'}>
                        <CustomIconMusicNoteEighth
                          aria-label={`default-icon-${song.title}`}
                          size={16}
                        />
                      </Center>
                    </Avatar>

                    <Stack gap={2}>
                      <Text fw={500} lineClamp={1} lh={'xxs'}>
                        {song.title}
                      </Text>
                      {song.artist && (
                        <Text fz={'xxs'} c={'dimmed'} fw={500} lineClamp={1} lh={'xxs'}>
                          {song.artist.name}
                        </Text>
                      )}
                    </Stack>
                  </Group>

                  {!songArrangementsMap.get(song.id) ? (
                    <Tooltip label={'This song will be skipped'}>
                      <Button
                        variant={'subtle'}
                        size={'compact-xs'}
                        disabled={true}
                        style={{ flexShrink: 0 }}
                      >
                        No Arrangement Found
                      </Button>
                    </Tooltip>
                  ) : (
                    <Menu>
                      <Menu.Target>
                        <Button variant={'subtle'} size={'compact-xs'} style={{ flexShrink: 0 }}>
                          {songArrangementsMap.get(song.id).name}
                        </Button>
                      </Menu.Target>

                      <Menu.Dropdown>
                        {song.arrangements.map((arrangement) => (
                          <Menu.Item
                            key={arrangement.id}
                            data-active={arrangement.id === songArrangementsMap.get(song.id).id}
                            sx={(theme) => ({
                              transition: '0.2s',
                              fontSize: theme.fontSizes.xs,
                              '.mantine-Menu-itemLabel': { display: 'flex', gap: 8 },
                              color: theme.colors.gray[7],
                              backgroundColor: alpha(theme.white, 0.7),
                              '.DefaultIcon': { color: theme.colors.gray[6] },
                              '&:hover': {
                                color: theme.colors.gray[8],
                                '.DefaultIcon': { color: alpha(theme.colors.gray[7], 0.7) },
                                backgroundColor: alpha(theme.colors.gray[1], 1)
                              },

                              ...(arrangement.id === songArrangementsMap.get(song.id).id && {
                                color: theme.colors.primary[4],
                                backgroundColor: alpha(theme.colors.primary[1], 0.4),
                                '.DefaultIcon': { color: theme.colors.primary[4] },
                                // same as above
                                '&:hover': {
                                  color: theme.colors.primary[4],
                                  '.DefaultIcon': { color: theme.colors.primary[4] },
                                  backgroundColor: alpha(theme.colors.primary[1], 0.4)
                                }
                              }),

                              ...(song.defaultArrangementId === arrangement.id && {
                                fontWeight: 500
                              })
                            })}
                            onClick={() => handleChangeArrangement(song.id, arrangement)}
                          >
                            {arrangement.name}
                            {song.defaultArrangementId === arrangement.id && (
                              <Center className={'DefaultIcon'}>
                                <IconStarFilled size={12} />
                              </Center>
                            )}
                          </Menu.Item>
                        ))}
                      </Menu.Dropdown>
                    </Menu>
                  )}
                </Group>
              ))}
          </Stack>
        </ScrollArea.Autosize>

        <Tooltip
          label={'To submit custom rehearsals, the songs need arrangements'}
          disabled={!noArrangementsFound}
        >
          <Button
            loading={isCustomLoading}
            disabled={noArrangementsFound}
            onClick={handleAddCustomRehearsals}
          >
            Submit
          </Button>
        </Tooltip>
      </Stack>
    </Modal>
  )
}

export default CustomRehearsalsModal
