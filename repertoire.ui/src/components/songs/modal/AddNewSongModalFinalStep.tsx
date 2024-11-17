import { Stack, TextInput } from '@mantine/core'
import { IconBrandYoutubeFilled, IconGuitarPickFilled } from '@tabler/icons-react'
import { FileWithPath } from '@mantine/dropzone'
import { UseFormReturnType } from '@mantine/form'
import { Dispatch, SetStateAction } from 'react'
import LargeImageDropzoneWithPreview from '../../image/LargeImageDropzoneWithPreview.tsx'
import { AddNewSongForm } from '../../../validation/songsForm.ts'

interface AddNewSongModalFinalStepProps {
  form: UseFormReturnType<AddNewSongForm, (values: AddNewSongForm) => AddNewSongForm>
  image: FileWithPath | null
  setImage: Dispatch<SetStateAction<FileWithPath | null>>
}

function AddNewSongModalFinalStep({ form, image, setImage }: AddNewSongModalFinalStepProps) {
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
    </Stack>
  )
}

export default AddNewSongModalFinalStep
