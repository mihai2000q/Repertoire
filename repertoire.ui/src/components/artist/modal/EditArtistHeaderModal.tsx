import Artist from '../../../types/models/Artist.ts'
import {
  Button,
  Checkbox,
  LoadingOverlay,
  Modal,
  Stack,
  Text,
  TextInput,
  Tooltip
} from '@mantine/core'
import {
  useDeleteImageFromArtistMutation,
  useSaveImageToArtistMutation,
  useUpdateArtistMutation
} from '../../../state/api/artistsApi.ts'
import { useEffect, useState } from 'react'
import { useForm } from '@mantine/form'
import { zod4Resolver } from 'mantine-form-zod-resolver'
import { EditArtistHeaderForm, editArtistHeaderSchema } from '../../../validation/artistsForm.ts'
import LargeImageDropzoneWithPreview from '../../@ui/image/LargeImageDropzoneWithPreview.tsx'
import { toast } from 'react-toastify'
import { useDidUpdate } from '@mantine/hooks'
import { FileWithPath } from '@mantine/dropzone'

interface EditArtistHeaderModalProps {
  artist: Artist
  opened: boolean
  onClose: () => void
}

function EditArtistHeaderModal({ artist, opened, onClose }: EditArtistHeaderModalProps) {
  const [updateArtistMutation, { isLoading: isUpdateLoading }] = useUpdateArtistMutation()
  const [saveImageMutation, { isLoading: isSaveImageLoading }] = useSaveImageToArtistMutation()
  const [deleteImageMutation, { isLoading: isDeleteImageLoading }] =
    useDeleteImageFromArtistMutation()
  const isLoading = isUpdateLoading || isSaveImageLoading || isDeleteImageLoading

  const [artistHasChanged, setArtistHasChanged] = useState(false)
  const [imageHasChanged, setImageHasChanged] = useState(false)
  const hasChanged = artistHasChanged || imageHasChanged

  const form = useForm<EditArtistHeaderForm>({
    mode: 'uncontrolled',
    initialValues: {
      name: artist.name,
      image: artist.imageUrl,
      isBand: artist.isBand
    },
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zod4Resolver(editArtistHeaderSchema),
    onValuesChange: (values) => {
      setArtistHasChanged(values.name !== artist.name || values.isBand !== artist.isBand)
      setImageHasChanged(values.image !== artist.imageUrl)
    }
  })

  const [image, setImage] = useState<string | FileWithPath>(artist.imageUrl)
  useEffect(() => form.setFieldValue('image', image), [image])

  useDidUpdate(() => {
    form.setValues({
      name: artist.name,
      image: artist.imageUrl,
      isBand: artist.isBand
    })
    setImage(artist.imageUrl)
  }, [artist])

  async function updateArtist({ name, image, isBand }) {
    if (artistHasChanged)
      await updateArtistMutation({
        id: artist.id,
        name: name.trim(),
        isBand: isBand
      }).unwrap()

    if (image !== null && typeof image !== 'string')
      await saveImageMutation({
        id: artist.id,
        image: image
      })
    else if (image === null && artist.imageUrl) await deleteImageMutation(artist.id)

    toast.info('Artist updated!')
    onClose()
    setArtistHasChanged(false)
    setImageHasChanged(false)
  }

  return (
    <Modal opened={opened} onClose={onClose} title={'Edit Artist Header'}>
      <LoadingOverlay visible={isLoading} loaderProps={{ type: 'bars' }} />

      <form onSubmit={form.onSubmit(updateArtist)}>
        <Stack px={'xs'} py={0}>
          <LargeImageDropzoneWithPreview
            image={image}
            setImage={setImage}
            defaultValue={artist.imageUrl}
          />

          <TextInput
            withAsterisk={true}
            maxLength={100}
            label="Name"
            placeholder="The name of the artist"
            key={form.key('name')}
            {...form.getInputProps('name')}
          />

          <Checkbox
            aria-label={'is-band'}
            label={
              <Text inline fz={'sm'}>
                The artist is a <b>band</b>
              </Text>
            }
            styles={{ label: { paddingLeft: 8 }, labelWrapper: { justifyContent: 'center' } }}
            key={form.key('isBand')}
            {...form.getInputProps('isBand', { type: 'checkbox' })}
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
    </Modal>
  )
}

export default EditArtistHeaderModal
