import Song from '../../../types/models/Song.ts'
import {
  Box,
  Button,
  Group,
  LoadingOverlay,
  Modal,
  Stack,
  Text,
  TextInput,
  Tooltip
} from '@mantine/core'
import {
  useDeleteImageFromSongMutation,
  useSaveImageToSongMutation,
  useUpdateSongMutation
} from '../../../state/songsApi.ts'
import { useEffect, useState } from 'react'
import { useForm, zodResolver } from '@mantine/form'
import { EditSongHeaderForm, editSongHeaderValidation } from '../../../validation/songsForm.ts'
import { DatePickerInput } from '@mantine/dates'
import { IconCalendarFilled, IconInfoCircleFilled } from '@tabler/icons-react'
import LargeImageDropzoneWithPreview from '../../@ui/image/LargeImageDropzoneWithPreview.tsx'
import { toast } from 'react-toastify'
import { FileWithPath } from '@mantine/dropzone'
import { useDidUpdate } from '@mantine/hooks'

interface EditSongHeaderModalProps {
  song: Song
  opened: boolean
  onClose: () => void
}

function EditSongHeaderModal({ song, opened, onClose }: EditSongHeaderModalProps) {
  const [updateSongMutation, { isLoading: isUpdateLoading }] = useUpdateSongMutation()
  const [saveImageMutation, { isLoading: isSaveImageLoading }] = useSaveImageToSongMutation()
  const [deleteImageMutation, { isLoading: isDeleteImageLoading }] =
    useDeleteImageFromSongMutation()
  const isLoading = isUpdateLoading || isSaveImageLoading || isDeleteImageLoading

  const [hasChanged, setHasChanged] = useState(false)

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      title: song.title,
      releaseDate: song.releaseDate && new Date(song.releaseDate),
      image: song.imageUrl
    } as EditSongHeaderForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(editSongHeaderValidation),
    onValuesChange: (values) => {
      setHasChanged(
        values.title !== song.title ||
          values.releaseDate?.toISOString() !== new Date(song.releaseDate).toISOString() ||
          values.image !== song.imageUrl
      )
    }
  })

  const [image, setImage] = useState<string | FileWithPath>(song.imageUrl)
  useEffect(() => form.setFieldValue('image', image), [image])
  useDidUpdate(() => setImage(song.imageUrl), [song])

  async function updateSong({ title, releaseDate, image }: EditSongHeaderForm) {
    title = title.trim()

    await updateSongMutation({
      ...song,
      guitarTuningId: song.guitarTuning?.id,
      id: song.id,
      title: title,
      releaseDate: releaseDate
    }).unwrap()

    if (image !== null && typeof image !== 'string') {
      await saveImageMutation({
        id: song.id,
        image: image
      })
    } else if (image === null && song.imageUrl) {
      await deleteImageMutation(song.id)
    }

    toast.info('Song header updated!')
    onClose()
    setHasChanged(false)
  }

  return (
    <Modal opened={opened} onClose={onClose} title={'Edit Song Header'}>
      <Modal.Body px={'xs'} py={0}>
        <LoadingOverlay visible={isLoading} />

        <form onSubmit={form.onSubmit(updateSong)}>
          <Stack>
            <LargeImageDropzoneWithPreview
              image={image}
              setImage={setImage}
              defaultValue={song.imageUrl}
            />

            {!image && song.album?.imageUrl && (
              <Group gap={6}>
                <Box c={'primary.8'} mt={3}>
                  <IconInfoCircleFilled size={15} />
                </Box>

                <Text inline fw={500} c={'dimmed'} fz={'xs'}>
                  The song image is inherited from the album
                </Text>
              </Group>
            )}

            <TextInput
              withAsterisk={true}
              maxLength={100}
              label="Title"
              placeholder="The title of the song"
              key={form.key('title')}
              {...form.getInputProps('title')}
            />

            <DatePickerInput
              label={'Release Date'}
              leftSection={<IconCalendarFilled size={20} />}
              placeholder={'Choose the release date'}
              key={form.key('releaseDate')}
              {...form.getInputProps('releaseDate')}
            />

            <Tooltip
              disabled={hasChanged}
              label={'You need to make a change before saving'}
              position="bottom"
            >
              <Button
                type={'submit'}
                data-disabled={!hasChanged}
                onClick={(e) => (!hasChanged ? e.preventDefault() : {})}
              >
                Save
              </Button>
            </Tooltip>
          </Stack>
        </form>
      </Modal.Body>
    </Modal>
  )
}

export default EditSongHeaderModal
