import { alpha, Button, Center, Group, Modal, Stack, Text, TextInput, Tooltip } from '@mantine/core'
import { useState } from 'react'
import { FileWithPath } from '@mantine/dropzone'
import { useForm, zodResolver } from '@mantine/form'
import { toast } from 'react-toastify'
import { useCreateSongMutation, useSaveImageToSongMutation } from '../../../state/api/songsApi.ts'
import ImageDropzoneWithPreview from '../../@ui/image/ImageDropzoneWithPreview.tsx'
import { AddNewAlbumSongForm, addNewAlbumSongValidation } from '../../../validation/albumsForm.ts'
import Album from '../../../types/models/Album.ts'
import {
  IconBrandYoutubeFilled,
  IconGuitarPickFilled,
  IconInfoCircleFilled,
  IconStarFilled
} from '@tabler/icons-react'
import GuitarTuningSelectButton from '../../@ui/form/select/GuitarTuningSelectButton.tsx'
import { GuitarTuning } from '../../../types/models/Song.ts'
import DifficultySelectButton from '../../@ui/form/select/DifficultySelectButton.tsx'
import Difficulty from '../../../types/enums/Difficulty.ts'
import CustomIconGuitarHead from '../../@ui/icons/CustomIconGuitarHead.tsx'
import NumberInputButton from '../../@ui/form/input/NumberInputButton.tsx'
import CustomIconMetronome from '../../@ui/icons/CustomIconMetronome.tsx'
import TextInputButton from '../../@ui/form/input/TextInputButton.tsx'

interface AddNewAlbumSongModalProps {
  opened: boolean
  onClose: () => void
  album: Album | undefined
}

function AddNewAlbumSongModal({ opened, onClose, album }: AddNewAlbumSongModalProps) {
  const [createSongMutation, { isLoading: isCreateSongLoading }] = useCreateSongMutation()
  const [saveImageMutation, { isLoading: isSaveImageLoading }] = useSaveImageToSongMutation()
  const isLoading = isCreateSongLoading || isSaveImageLoading

  const [image, setImage] = useState<FileWithPath>(null)
  const [guitarTuning, setGuitarTuning] = useState<GuitarTuning>(null)
  const [difficulty, setDifficulty] = useState<Difficulty>(null)
  const [bpm, setBpm] = useState<string | number>(null)

  const [isSongsterrLinkSelected, setIsSongsterrLinkSelected] = useState(false)
  const [isYoutubeLinkSelected, setIsYoutubeLinkSelected] = useState(false)

  const onCloseWithImage = () => {
    onClose()
    setImage(null)
  }

  const inheritedValues = [
    ...(album?.releaseDate ? ['release date'] : []),
    ...(album?.artist ? ['artist'] : [])
  ]

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      title: ''
    } as AddNewAlbumSongForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(addNewAlbumSongValidation),
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

  async function addSong({ title, songsterrLink, youtubeLink }: AddNewAlbumSongForm) {
    title = title.trim()
    songsterrLink = songsterrLink?.trim() === '' ? undefined : songsterrLink?.trim()
    youtubeLink = youtubeLink?.trim() === '' ? undefined : youtubeLink?.trim()

    const res = await createSongMutation({
      title: title,
      description: '',
      albumId: album?.id,
      guitarTuningId: guitarTuning?.id,
      difficulty: difficulty,
      bpm: typeof bpm === 'string' ? undefined : bpm,
      songsterrLink: songsterrLink,
      youtubeLink: youtubeLink
    }).unwrap()

    if (image) await saveImageMutation({ image: image, id: res.id }).unwrap()

    toast.success(`${title} added!`)

    onCloseWithImage()
    form.reset()
    setGuitarTuning(null)
    setDifficulty(null)
    setBpm(null)
    setIsSongsterrLinkSelected(false)
    setIsYoutubeLinkSelected(false)
  }

  return (
    <Modal opened={opened} onClose={onCloseWithImage} title={'Add New Song'}>
      <Modal.Body p={'xs'}>
        <form onSubmit={form.onSubmit(addSong)}>
          <Stack>
            <Group align={'start'}>
              <ImageDropzoneWithPreview image={image} setImage={setImage} />

              <Stack flex={1} gap={'xxs'}>
                <TextInput
                  withAsterisk={true}
                  maxLength={100}
                  label="Title"
                  placeholder="The title of the song"
                  key={form.key('title')}
                  {...form.getInputProps('title')}
                />

                <Group gap={2}>
                  <Tooltip.Group openDelay={500}>
                    <GuitarTuningSelectButton
                      guitarTuning={guitarTuning}
                      setGuitarTuning={setGuitarTuning}
                      size={'md'}
                      icon={<CustomIconGuitarHead size={16} />}
                    />
                    <DifficultySelectButton
                      difficulty={difficulty}
                      setDifficulty={setDifficulty}
                      size={'md'}
                      icon={<IconStarFilled size={16} />}
                    />
                    <NumberInputButton
                      icon={<CustomIconMetronome size={16} />}
                      aria-label={'bpm'}
                      size={'md'}
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
                      size={'md'}
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
                      size={'md'}
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
                      <Center c={'primary.8'} mb={1}>
                        <IconInfoCircleFilled size={13} />
                      </Center>

                      <Text fw={500} c={'dimmed'} fz={'xs'}>
                        If no image is uploaded, it will be inherited.
                      </Text>
                    </Group>
                  )}
                  {inheritedValues.length > 0 && (
                    <Group gap={'xxs'} wrap={'nowrap'}>
                      <Center c={'primary.8'}>
                        <IconInfoCircleFilled size={13} />
                      </Center>

                      <Text fw={500} c={'dimmed'} fz={'xs'}>
                        The new song will inherit the <b>{inheritedValues.join(', ')}</b>.
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
      </Modal.Body>
    </Modal>
  )
}

export default AddNewAlbumSongModal
