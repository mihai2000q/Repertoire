import { Button, Group, Modal, Stack, TextInput } from '@mantine/core'
import { useState } from 'react'
import { FileWithPath } from '@mantine/dropzone'
import { useForm, zodResolver } from '@mantine/form'
import { toast } from 'react-toastify'
import { AddNewAlbumForm, addNewAlbumValidation } from '../../../validation/albumsForm.ts'
import {
  useCreateAlbumMutation,
  useSaveImageToAlbumMutation
} from '../../../state/api/albumsApi.ts'
import LargeImageDropzoneWithPreview from '../../@ui/image/LargeImageDropzoneWithPreview.tsx'
import { DatePickerInput } from '@mantine/dates'
import ArtistAutocomplete from '../../@ui/form/input/ArtistAutocomplete.tsx'
import { IconCalendarRepeat } from '@tabler/icons-react'
import { ArtistSearch } from '../../../types/models/Search.ts'

interface AddNewAlbumModalProps {
  opened: boolean
  onClose: () => void
}

function AddNewAlbumModal({ opened, onClose }: AddNewAlbumModalProps) {
  const [createAlbumMutation, { isLoading: isCreateAlbumLoading }] = useCreateAlbumMutation()
  const [saveImageMutation, { isLoading: isSaveImageLoading }] = useSaveImageToAlbumMutation()
  const isLoading = isCreateAlbumLoading || isSaveImageLoading

  const [artist, setArtist] = useState<ArtistSearch>(null)

  const [image, setImage] = useState<FileWithPath>(null)

  const onCloseWithImage = () => {
    onClose()
    setImage(null)
  }

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      title: ''
    } as AddNewAlbumForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(addNewAlbumValidation)
  })

  async function addAlbum({ title, artistName, releaseDate }: AddNewAlbumForm) {
    title = title.trim()
    artistName = artistName?.trim() === '' ? null : artistName?.trim()

    const res = await createAlbumMutation({
      title: title,
      releaseDate: releaseDate,
      artistId: artist?.id,
      artistName: artist ? undefined : artistName
    }).unwrap()

    if (image) await saveImageMutation({ image: image, id: res.id }).unwrap()

    toast.success(`${title} added!`)

    onCloseWithImage()
    form.reset()
    setArtist(null)
  }

  return (
    <Modal opened={opened} onClose={onCloseWithImage} title={'Add New Album'}>
      <Modal.Body p={'xs'}>
        <form onSubmit={form.onSubmit(addAlbum)}>
          <Stack>
            <TextInput
              withAsterisk={true}
              maxLength={100}
              label="Title"
              placeholder="The title of the album"
              key={form.key('title')}
              {...form.getInputProps('title')}
            />

            <Group>
              <ArtistAutocomplete
                artist={artist}
                setArtist={setArtist}
                key={form.key('artistName')}
                setValue={(v) => form.setFieldValue('artistName', v)}
                {...form.getInputProps('artistName')}
              />

              <DatePickerInput
                flex={1}
                label={'Release Date'}
                leftSection={<IconCalendarRepeat size={20} />}
                placeholder={'Choose the release date'}
                key={form.key('releaseDate')}
                {...form.getInputProps('releaseDate')}
              />
            </Group>

            <LargeImageDropzoneWithPreview image={image} setImage={setImage} />

            <Button style={{ alignSelf: 'end' }} type={'submit'} loading={isLoading}>
              Submit
            </Button>
          </Stack>
        </form>
      </Modal.Body>
    </Modal>
  )
}

export default AddNewAlbumModal
