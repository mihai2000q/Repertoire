import { alpha, Button, Center, Group, Modal, Stack, Text, TextInput, Tooltip } from '@mantine/core'
import { useState } from 'react'
import { FileWithPath } from '@mantine/dropzone'
import { useForm } from '@mantine/form'
import { zod4Resolver } from 'mantine-form-zod-resolver'
import { AddNewArtistSongForm, addNewArtistSongSchema } from '../../../validation/artistsForm.ts'
import { toast } from 'react-toastify'
import { useCreateSongMutation, useSaveImageToSongMutation } from '../../../state/api/songsApi.ts'
import ImageDropzoneWithPreview from '../../@ui/image/ImageDropzoneWithPreview.tsx'
import GuitarTuningSelectButton from '../../@ui/form/select/button/GuitarTuningSelectButton.tsx'
import DifficultySelectButton from '../../@ui/form/select/button/DifficultySelectButton.tsx'
import {
  IconBrandYoutubeFilled,
  IconCalendarRepeat,
  IconDisc,
  IconGuitarPickFilled,
  IconInfoCircleFilled
} from '@tabler/icons-react'
import NumberInputButton from '../../@ui/form/input/button/NumberInputButton.tsx'
import CustomIconMetronome from '../../@ui/icons/CustomIconMetronome.tsx'
import TextInputButton from '../../@ui/form/input/button/TextInputButton.tsx'
import { GuitarTuning } from '../../../types/models/Song.ts'
import Difficulty from '../../../types/enums/Difficulty.ts'
import DatePickerButton from '../../@ui/form/date/DatePickerButton.tsx'
import dayjs from 'dayjs'
import { AlbumSearch } from '../../../types/models/Search.ts'
import AlbumSelectButton from '../../@ui/form/select/button/AlbumSelectButton.tsx'

interface AddNewArtistSongModalProps {
  opened: boolean
  onClose: () => void
  artistId: string | undefined
}

function AddNewArtistSongModal({ opened, onClose, artistId }: AddNewArtistSongModalProps) {
  const [createSongMutation, { isLoading: isCreateSongLoading }] = useCreateSongMutation()
  const [saveImageMutation, { isLoading: isSaveImageLoading }] = useSaveImageToSongMutation()
  const isLoading = isCreateSongLoading || isSaveImageLoading

  const [image, setImage] = useState<FileWithPath>(null)
  const [album, setAlbum] = useState<AlbumSearch>(null)
  const [releaseDate, setReleaseDate] = useState<string>()
  const [guitarTuning, setGuitarTuning] = useState<GuitarTuning>(null)
  const [difficulty, setDifficulty] = useState<Difficulty>(null)
  const [bpm, setBpm] = useState<string | number>()

  const [isSongsterrLinkSelected, setIsSongsterrLinkSelected] = useState(false)
  const [isYoutubeLinkSelected, setIsYoutubeLinkSelected] = useState(false)

  const onCloseWithImage = () => {
    onClose()
    setImage(null)
  }

  const form = useForm<AddNewArtistSongForm>({
    mode: 'uncontrolled',
    initialValues: {
      title: ''
    },
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zod4Resolver(addNewArtistSongSchema),
    onValuesChange: (values) => {
      setIsSongsterrLinkSelected(
        values.songsterrLink !== '' &&
          values.songsterrLink !== null &&
          values.songsterrLink !== undefined
      )
      setIsYoutubeLinkSelected(
        values.youtubeLink !== '' && values.youtubeLink !== null && values.youtubeLink !== undefined
      )
    }
  })

  async function addSong({ title, youtubeLink, songsterrLink }: AddNewArtistSongForm) {
    title = title.trim()
    songsterrLink = songsterrLink?.trim() === '' ? undefined : songsterrLink?.trim()
    youtubeLink = youtubeLink?.trim() === '' ? undefined : youtubeLink?.trim()

    const res = await createSongMutation({
      title: title,
      description: '',
      artistId: album ? undefined : artistId,
      albumId: album?.id,
      releaseDate: releaseDate,
      guitarTuningId: guitarTuning?.id,
      difficulty: difficulty ?? undefined,
      bpm: typeof bpm === 'number' ? bpm : undefined,
      songsterrLink: songsterrLink,
      youtubeLink: youtubeLink
    }).unwrap()

    if (image) await saveImageMutation({ image: image, id: res.id }).unwrap()

    toast.success(`${title} added!`)
    onCloseWithImage()
    form.reset()
    setAlbum(null)
    setReleaseDate(null)
    setGuitarTuning(null)
    setDifficulty(null)
    setBpm(null)
    setIsSongsterrLinkSelected(false)
    setIsYoutubeLinkSelected(false)
  }

  return (
    <Modal opened={opened} onClose={onCloseWithImage} title={'Add New Song'}>
      <form onSubmit={form.onSubmit(addSong)}>
        <Stack p={'xs'}>
          <Group align={'start'}>
            <ImageDropzoneWithPreview image={image} setImage={setImage} />

            <Stack flex={1} gap={'xxs'}>
              <Group gap={'xxs'}>
                <TextInput
                  flex={1}
                  withAsterisk={true}
                  maxLength={100}
                  label="Title"
                  placeholder="The title of the song"
                  key={form.key('title')}
                  {...form.getInputProps('title')}
                />
                <AlbumSelectButton
                  mt={form.getInputProps('title').error ? 3 : 19}
                  size={'lg'}
                  icon={<IconDisc size={20} />}
                  album={album}
                  setAlbum={setAlbum}
                  searchFilter={artistId ? [`artist.id = ${artistId}`] : []}
                />
              </Group>

              <Group gap={2}>
                <Tooltip.Group openDelay={500}>
                  <DatePickerButton
                    aria-label={'release-date'}
                    icon={<IconCalendarRepeat size={18} />}
                    value={releaseDate}
                    onChange={setReleaseDate}
                    tooltipLabels={{
                      default: 'Select a release date',
                      selected: (val) => `Released on ${dayjs(val).format('D MMMM YYYY')}`
                    }}
                  />
                  <GuitarTuningSelectButton
                    guitarTuning={guitarTuning}
                    setGuitarTuning={setGuitarTuning}
                  />
                  <DifficultySelectButton difficulty={difficulty} setDifficulty={setDifficulty} />
                  <NumberInputButton
                    icon={<CustomIconMetronome size={16} />}
                    aria-label={'bpm'}
                    inputProps={{
                      'aria-label': 'bpm',
                      placeholder: 'Enter bpm',
                      leftSection: <CustomIconMetronome size={15} />,
                      value: bpm,
                      onChange: setBpm
                    }}
                    tooltipLabels={{
                      selected: `Bpm is ${bpm}`,
                      default: 'Enter a bpm'
                    }}
                  />
                  <TextInputButton
                    icon={<IconGuitarPickFilled size={16} />}
                    inputKey={form.key('songsterrLink')}
                    aria-label={'songsterr'}
                    inputProps={{
                      'aria-label': 'songsterr',
                      placeholder: 'Enter Songsterr Link',
                      leftSection: <IconGuitarPickFilled size={15} />,
                      w: 370,
                      ...form.getInputProps('songsterrLink')
                    }}
                    isSelected={isSongsterrLinkSelected}
                    tooltipLabels={{
                      default: 'Enter Songsterr Link',
                      selected: 'Songsterr Link entered!'
                    }}
                    variant={'transparent'}
                    sx={(theme) => ({
                      ...(Object.entries(theme.components.ActionIcon.styles(theme).root).find(
                        (s) => s[0] === '&[data-variant="form"]'
                      )[1] as object),

                      '&[aria-selected="true"]': {
                        color: theme.colors.blue[5],
                        backgroundColor: alpha(theme.colors.blue[1], 0.5),

                        '&:hover': {
                          color: theme.colors.blue[6],
                          backgroundColor: theme.colors.blue[1]
                        }
                      }
                    })}
                  />
                  <TextInputButton
                    icon={<IconBrandYoutubeFilled size={16} />}
                    inputKey={form.key('youtubeLink')}
                    aria-label={'youtube'}
                    inputProps={{
                      'aria-label': 'youtube',
                      placeholder: 'Enter Youtube Link',
                      leftSection: <IconBrandYoutubeFilled size={15} />,
                      w: 300,
                      ...form.getInputProps('youtubeLink')
                    }}
                    isSelected={isYoutubeLinkSelected}
                    tooltipLabels={{
                      default: 'Enter Youtube Link',
                      selected: 'Youtube Link entered!'
                    }}
                    variant={'transparent'}
                    sx={(theme) => ({
                      ...(Object.entries(theme.components.ActionIcon.styles(theme).root).find(
                        (s) => s[0] === '&[data-variant="form"]'
                      )[1] as object),
                      '&[aria-selected="true"]': {
                        color: theme.colors.red[5],
                        backgroundColor: alpha(theme.colors.red[1], 0.5),

                        '&:hover': {
                          color: theme.colors.red[6],
                          backgroundColor: theme.colors.red[1]
                        }
                      }
                    })}
                  />
                </Tooltip.Group>
              </Group>

              <Stack gap={0}>
                {!image && album?.imageUrl && (
                  <Group gap={'xxs'}>
                    <Center c={'primary.8'}>
                      <IconInfoCircleFilled size={13} />
                    </Center>

                    <Text fw={500} c={'dimmed'} fz={'xs'}>
                      The image will be inherited from album.
                    </Text>
                  </Group>
                )}
                {!releaseDate && album?.releaseDate && (
                  <Group gap={'xxs'}>
                    <Center c={'primary.8'}>
                      <IconInfoCircleFilled size={13} />
                    </Center>

                    <Text fw={500} c={'dimmed'} fz={'xs'}>
                      The release date will be inherited from album.
                    </Text>
                  </Group>
                )}
              </Stack>
            </Stack>
          </Group>

          <Button style={{ alignSelf: 'center' }} type={'submit'} loading={isLoading}>
            Submit
          </Button>
        </Stack>
      </form>
    </Modal>
  )
}

export default AddNewArtistSongModal
