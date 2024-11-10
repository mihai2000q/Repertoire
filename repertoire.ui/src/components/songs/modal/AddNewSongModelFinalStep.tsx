import {
  ActionIcon,
  FileButton,
  Group,
  Image,
  Stack,
  Text,
  TextInput,
  Tooltip
} from '@mantine/core'
import {
  IconBrandYoutubeFilled,
  IconGuitarPickFilled,
  IconPhoto,
  IconPhotoDown,
  IconUpload,
  IconX
} from '@tabler/icons-react'
import { Dropzone, FileWithPath, IMAGE_MIME_TYPE } from '@mantine/dropzone'
import { UseFormReturnType } from '@mantine/form'
import { Dispatch, SetStateAction } from 'react'

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

      {image ? (
        <Group justify={'center'} align={'center'}>
          <Image
            src={URL.createObjectURL(image)}
            h={200}
            w={320}
            radius={'md'}
            alt={'song-image'}
          />

          <FileButton onChange={setImage} accept="image/png,image/jpeg">
            {(props) => (
              <Tooltip label={'Reload Image'} openDelay={300} position={'right'}>
                <ActionIcon
                  c={'dark'}
                  variant={'subtle'}
                  aria-label={'add-image-button'}
                  size={'xl'}
                  {...props}
                >
                  <IconPhotoDown size={20} />
                </ActionIcon>
              </Tooltip>
            )}
          </FileButton>
        </Group>
      ) : (
        <Dropzone
          onDrop={(files) => setImage(files[0])}
          accept={IMAGE_MIME_TYPE}
          multiple={false}
          sx={(theme) => ({
            cursor: 'pointer',
            borderRadius: '16px',
            transition: '0.3s',
            color: theme.colors.gray[6],
            '&:hover': {
              color: theme.colors.gray[7],
              backgroundColor: theme.colors.gray[2]
            },

            '&:where([data-accept])': {
              color: theme.colors.green[6],
              backgroundColor: theme.colors.green[1]
            },

            '&:where([data-reject])': {
              color: theme.colors.red[6],
              backgroundColor: theme.colors.red[1]
            }
          })}
        >
          <Group justify="center" gap="xs" h={150} style={{ pointerEvents: 'none' }}>
            <Dropzone.Accept>
              <IconUpload size={40} />
            </Dropzone.Accept>
            <Dropzone.Reject>
              <IconX size={40} />
            </Dropzone.Reject>
            <Dropzone.Idle>
              <IconPhoto size={40} />
            </Dropzone.Idle>

            <Text size="xl" inline>
              Drag image here or click to select it
            </Text>
          </Group>
        </Dropzone>
      )}
    </Stack>
  )
}

export default AddNewSongModelFinalStep
