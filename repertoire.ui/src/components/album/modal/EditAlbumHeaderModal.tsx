import Album from '../../../types/models/Album.ts'
import {
  Button,
  Center,
  Group,
  LoadingOverlay,
  Modal,
  Stack,
  Text,
  TextInput,
  Tooltip
} from '@mantine/core'
import {
  useDeleteImageFromAlbumMutation,
  useSaveImageToAlbumMutation,
  useUpdateAlbumMutation
} from '../../../state/api/albumsApi.ts'
import { useEffect, useState } from 'react'
import { useForm } from '@mantine/form'
import { zod4Resolver } from 'mantine-form-zod-resolver'
import { EditAlbumHeaderForm, editAlbumHeaderSchema } from '../../../validation/albumsForm.ts'
import { DatePickerInput } from '@mantine/dates'
import { IconCalendarRepeat, IconInfoCircleFilled } from '@tabler/icons-react'
import LargeImageDropzoneWithPreview from '../../@ui/image/LargeImageDropzoneWithPreview.tsx'
import { toast } from 'react-toastify'
import { useDidUpdate } from '@mantine/hooks'
import { FileWithPath } from '@mantine/dropzone'
import ArtistSelect from '../../@ui/form/select/ArtistSelect.tsx'
import { ArtistSearch } from '../../../types/models/Search.ts'

interface EditAlbumHeaderModalProps {
  album: Album
  opened: boolean
  onClose: () => void
}

function EditAlbumHeaderModal({ album, opened, onClose }: EditAlbumHeaderModalProps) {
  const [updateAlbumMutation, { isLoading: isUpdateLoading }] = useUpdateAlbumMutation()
  const [saveImageMutation, { isLoading: isSaveImageLoading }] = useSaveImageToAlbumMutation()
  const [deleteImageMutation, { isLoading: isDeleteImageLoading }] =
    useDeleteImageFromAlbumMutation()
  const isLoading = isUpdateLoading || isSaveImageLoading || isDeleteImageLoading

  const [albumHasChanged, setAlbumHasChanged] = useState(false)
  const [imageHasChanged, setImageHasChanged] = useState(false)
  const hasChanged = albumHasChanged || imageHasChanged

  const form = useForm<EditAlbumHeaderForm>({
    mode: 'uncontrolled',
    initialValues: {
      title: album.title,
      releaseDate: album.releaseDate,
      image: album.imageUrl,
      artistId: album.artist?.id
    },
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zod4Resolver(editAlbumHeaderSchema),
    onValuesChange: (values) => {
      setAlbumHasChanged(
        values.title !== album.title ||
          values.releaseDate !== album.releaseDate ||
          values.artistId !== album.artist?.id
      )
      setImageHasChanged(values.image !== album.imageUrl)
    }
  })

  const [image, setImage] = useState<string | FileWithPath>(album.imageUrl)
  useEffect(() => form.setFieldValue('image', image), [image])

  const [artist, setArtist] = useState(album.artist as unknown as ArtistSearch)
  useDidUpdate(() => form.setFieldValue('artistId', artist?.id), [artist])

  useDidUpdate(() => {
    form.setValues({
      title: album.title,
      releaseDate: album.releaseDate,
      image: album.imageUrl,
      artistId: album.artist?.id
    })
    setArtist(album.artist as unknown as ArtistSearch)
    setImage(album.imageUrl)
  }, [album])

  async function updateAlbum({ title, releaseDate, image, artistId }: EditAlbumHeaderForm) {
    if (albumHasChanged)
      await updateAlbumMutation({
        id: album.id,
        title: title.trim(),
        releaseDate: releaseDate,
        artistId: artistId
      }).unwrap()

    if (image && typeof image !== 'string')
      await saveImageMutation({
        id: album.id,
        image: image as FileWithPath
      })
    else if (!image && album.imageUrl) await deleteImageMutation(album.id)

    toast.info('Album updated!')
    onClose()
    setAlbumHasChanged(false)
    setImageHasChanged(false)
  }

  return (
    <Modal opened={opened} onClose={onClose} title={'Edit Album Header'}>
      <LoadingOverlay visible={isLoading} loaderProps={{ type: 'bars' }} />

      <form onSubmit={form.onSubmit(updateAlbum)}>
        <Stack px={'xs'} py={0}>
          <LargeImageDropzoneWithPreview
            image={image}
            setImage={setImage}
            defaultValue={album.imageUrl}
          />

          {album.imageUrl !== image && (
            <Group gap={'xxs'}>
              <Center c={'primary.8'}>
                <IconInfoCircleFilled size={13} />
              </Center>

              <Text inline fw={500} c={'dimmed'} fz={'xs'}>
                This change will update all the associated songs
              </Text>
            </Group>
          )}

          <TextInput
            withAsterisk={true}
            maxLength={100}
            label="Title"
            placeholder="The title of the album"
            key={form.key('title')}
            {...form.getInputProps('title')}
          />

          <Stack gap={6}>
            <Group gap={'sm'}>
              <ArtistSelect flex={1} artist={artist} setArtist={setArtist} />

              <DatePickerInput
                flex={1}
                label={'Release Date'}
                leftSection={<IconCalendarRepeat size={20} />}
                placeholder={'Choose the release date'}
                key={form.key('releaseDate')}
                {...form.getInputProps('releaseDate')}
              />
            </Group>

            {album.artist?.id !== artist?.id && (
              <Group gap={'xxs'}>
                <Center c={'primary.8'}>
                  <IconInfoCircleFilled size={13} />
                </Center>

                <Text inline fw={500} c={'dimmed'} fz={'xs'}>
                  This change will update all the associated songs&#39; artist
                </Text>
              </Group>
            )}
          </Stack>

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

export default EditAlbumHeaderModal
