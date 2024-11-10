import { Stack, TextInput } from '@mantine/core'
import { IconBrandYoutubeFilled, IconGuitarPickFilled } from '@tabler/icons-react'
import { FileWithPath } from '@mantine/dropzone'
import { UseFormReturnType } from '@mantine/form'
import { Dispatch, SetStateAction } from 'react'
import ImageDropzoneWithPreview from '../../image/ImageDropzoneWithPreview.tsx'

interface AddNewSongModelFinalStepProps {
  form: UseFormReturnType<unknown, (values: unknown) => unknown>
  image: FileWithPath | null
  setImage: Dispatch<SetStateAction<FileWithPath | null>>
}

function AddNewSongModelFinalStep({ form, image, setImage }: AddNewSongModelFinalStepProps) {
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

      <ImageDropzoneWithPreview image={image} setImage={setImage} />
    </Stack>
  )
}

export default AddNewSongModelFinalStep
