import { Dispatch, SetStateAction } from 'react'
import { ActionIcon, alpha, FileButton, Group, Image, Stack, Text, Tooltip } from '@mantine/core'
import { IconPhoto, IconPhotoDown, IconUpload, IconX } from '@tabler/icons-react'
import { Dropzone, FileWithPath, IMAGE_MIME_TYPE } from '@mantine/dropzone'

interface ImageDropzoneWithPreviewProps {
  image: FileWithPath
  setImage: Dispatch<SetStateAction<FileWithPath>>
}

function LargeImageDropzoneWithPreview({ image, setImage }: ImageDropzoneWithPreviewProps) {
  if (image) {
    return (
      <Group justify={'center'} align={'center'}>
        <Image
          src={URL.createObjectURL(image)}
          h={200}
          w={320}
          radius={'md'}
          alt={'image-preview'}
        />

        <Stack>
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

          <Tooltip label={'Remove Image'} openDelay={300} position={'bottom'}>
            <ActionIcon
              variant={'subtle'}
              color={'red.6'}
              aria-label={'remove-image-button'}
              size={'xl'}
              sx={(theme) => ({
                transition: '0.15s',
                '&:hover': { backgroundColor: alpha(theme.colors.red[5], 0.2) }
              })}
              onClick={() => setImage(null)}
            >
              <IconX size={18} />
            </ActionIcon>
          </Tooltip>
        </Stack>
      </Group>
    )
  }

  return (
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
  )
}

export default LargeImageDropzoneWithPreview
