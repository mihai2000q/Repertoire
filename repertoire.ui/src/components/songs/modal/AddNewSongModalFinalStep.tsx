import { Box, Group, Stack, Text, TextInput } from '@mantine/core'
import {
  IconBrandYoutubeFilled,
  IconGuitarPickFilled,
  IconInfoCircleFilled
} from '@tabler/icons-react'
import { FileWithPath } from '@mantine/dropzone'
import { UseFormReturnType } from '@mantine/form'
import { Dispatch, SetStateAction } from 'react'
import LargeImageDropzoneWithPreview from '../../@ui/image/LargeImageDropzoneWithPreview.tsx'
import { AddNewSongForm } from '../../../validation/songsForm.ts'
import Album from '../../../types/models/Album.ts'

interface AddNewSongModalFinalStepProps {
  form: UseFormReturnType<AddNewSongForm>
  image: FileWithPath | null
  setImage: Dispatch<SetStateAction<FileWithPath | null>>
  album: Album | null
}

function AddNewSongModalFinalStep({ form, image, setImage, album }: AddNewSongModalFinalStepProps) {
  return (
    <Stack>
      <TextInput
        leftSection={<IconGuitarPickFilled size={20} />}
        label="Songsterr"
        placeholder="Songsterr link"
        key={form.key('songsterrLink')}
        {...form.getInputProps('songsterrLink')}
      />
      <TextInput
        leftSection={<IconBrandYoutubeFilled size={20} />}
        label="Youtube"
        placeholder="Youtube link"
        key={form.key('youtubeLink')}
        {...form.getInputProps('youtubeLink')}
      />

      <LargeImageDropzoneWithPreview image={image} setImage={setImage} />

      {!image && album?.imageUrl && (
        <Group gap={6}>
          <Box c={'primary.8'} mt={3}>
            <IconInfoCircleFilled size={15} />
          </Box>

          <Text inline fw={500} c={'dimmed'} fz={'xs'}>
            If no image is uploaded, it will be inherited from the album
          </Text>
        </Group>
      )}
    </Stack>
  )
}

export default AddNewSongModalFinalStep
