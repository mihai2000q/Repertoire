import { Autocomplete, Button, Group, Loader, Modal, Stack, Text, TextInput } from '@mantine/core'
import { useState } from 'react'
import { FileWithPath } from '@mantine/dropzone'
import { useForm, zodResolver } from '@mantine/form'
import { toast } from 'react-toastify'
import { AddNewAlbumForm, addNewAlbumValidation } from '../../../validation/albumsForm.ts'
import { useCreateAlbumMutation, useSaveImageToAlbumMutation } from '../../../state/albumsApi.ts'
import { useGetArtistsQuery } from '../../../state/artistsApi.ts'
import ImageDropzoneWithPreview from '../../ImageDropzoneWithPreview.tsx'
import { DatePickerInput } from '@mantine/dates'

interface AddNewAlbumModalProps {
  opened: boolean
  onClose: () => void
}

function AddNewAlbumModal({ opened, onClose }: AddNewAlbumModalProps) {
  const [createAlbumMutation, { isLoading: isCreateAlbumLoading }] = useCreateAlbumMutation()
  const [saveImageMutation, { isLoading: isSaveImageLoading }] = useSaveImageToAlbumMutation()
  const isLoading = isCreateAlbumLoading || isSaveImageLoading

  const { data: artistsData, isLoading: isArtistsLoading } = useGetArtistsQuery({})
  const artists = artistsData?.models?.map((artist) => ({
    value: artist.id,
    label: artist.name
  }))

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

  async function addAlbum({ title, releaseDate }: AddNewAlbumForm) {
    title = title.trim()

    const res = await createAlbumMutation({ title, releaseDate }).unwrap()

    if (image) await saveImageMutation({ image: image, id: res.id }).unwrap()

    toast.success(`${title} has been added!`)

    onCloseWithImage()
    form.reset()
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
              <DatePickerInput
                flex={1}
                label={'Release Date'}
                placeholder={'Choose the release date'}
                key={form.key('releaseDate')}
                {...form.getInputProps('releaseDate')}
              />

              {isArtistsLoading ? (
                <Group gap={'xs'} flex={1}>
                  <Loader size={25} />
                  <Text fz={'sm'} c={'dimmed'}>
                    Loading Artists...
                  </Text>
                </Group>
              ) : (
                <Autocomplete
                  flex={1}
                  data={artists}
                  label={'Artist'}
                  placeholder={`${artists.length > 0 ? 'Choose or Create Artist' : 'Enter New Artist Name'}`}
                />
              )}
            </Group>

            <ImageDropzoneWithPreview image={image} setImage={setImage} />

            <Button style={{ alignSelf: 'end' }} type={'submit'} disabled={isLoading}>
              Submit
            </Button>
          </Stack>
        </form>
      </Modal.Body>
    </Modal>
  )
}

export default AddNewAlbumModal
