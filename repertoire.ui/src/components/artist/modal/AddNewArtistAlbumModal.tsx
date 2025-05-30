import { Button, Group, Modal, Stack, TextInput } from '@mantine/core'
import { useState } from 'react'
import { FileWithPath } from '@mantine/dropzone'
import { useForm } from '@mantine/form'
import { zod4Resolver } from 'mantine-form-zod-resolver'
import {
  AddNewArtistAlbumForm,
  addNewArtistAlbumSchema
} from '../../../validation/artistsForm.ts'
import {
  useCreateAlbumMutation,
  useSaveImageToAlbumMutation
} from '../../../state/api/albumsApi.ts'
import ImageDropzoneWithPreview from '../../@ui/image/ImageDropzoneWithPreview.tsx'
import { IconCalendarCheck, IconCalendarRepeat, IconDisc } from '@tabler/icons-react'
import DatePickerButton from '../../@ui/form/date/DatePickerButton.tsx'
import { toast } from 'react-toastify'
import dayjs from 'dayjs'

interface AddNewArtistAlbumModalProps {
  opened: boolean
  onClose: () => void
  artistId: string | undefined
}

function AddNewArtistAlbumModal({ opened, onClose, artistId }: AddNewArtistAlbumModalProps) {
  const [createAlbumMutation, { isLoading: isCreateAlbumLoading }] = useCreateAlbumMutation()
  const [saveImageMutation, { isLoading: isSaveImageLoading }] = useSaveImageToAlbumMutation()
  const isLoading = isCreateAlbumLoading || isSaveImageLoading

  const [image, setImage] = useState<FileWithPath>(null)
  const [releaseDate, setReleaseDate] = useState<string>()

  const onCloseWithImage = () => {
    onClose()
    setImage(null)
  }

  const form = useForm<AddNewArtistAlbumForm>({
    mode: 'uncontrolled',
    initialValues: {
      title: ''
    },
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zod4Resolver(addNewArtistAlbumSchema)
  })

  async function addAlbum({ title }) {
    title = title.trim()

    const res = await createAlbumMutation({
      title,
      releaseDate,
      artistId
    }).unwrap()

    if (image) await saveImageMutation({ image: image, id: res.id }).unwrap()

    toast.success(`${title} added!`)
    onCloseWithImage()
    form.reset()
    setReleaseDate(null)
  }

  return (
    <Modal opened={opened} onClose={onCloseWithImage} title={'Add New Album'}>
      <Modal.Body p={'xs'}>
        <form onSubmit={form.onSubmit(addAlbum)}>
          <Stack>
            <Group>
              <ImageDropzoneWithPreview
                image={image}
                setImage={setImage}
                icon={<IconDisc size={45} />}
              />
              <Group gap={'xxs'} flex={1}>
                <TextInput
                  flex={1}
                  withAsterisk={true}
                  maxLength={100}
                  label="Title"
                  placeholder="The title of the album"
                  key={form.key('title')}
                  {...form.getInputProps('title')}
                />
                <DatePickerButton
                  mt={form.getInputProps('title').error ? 3 : 19}
                  aria-label={'release-date'}
                  size={'lg'}
                  icon={<IconCalendarRepeat size={20} />}
                  successIcon={<IconCalendarCheck size={20} />}
                  value={releaseDate}
                  onChange={setReleaseDate}
                  tooltipLabels={{
                    default: 'Select a release date',
                    selected: (val) => `Released on ${dayjs(val).format('D MMMM YYYY')}`
                  }}
                />
              </Group>
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

export default AddNewArtistAlbumModal
