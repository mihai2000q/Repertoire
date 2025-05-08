import { Dispatch, SetStateAction } from 'react'
import {
  ActionIcon,
  alpha,
  AspectRatio,
  FileButton,
  Group,
  Image,
  Stack,
  Text,
  Tooltip
} from '@mantine/core'
import { IconPhoto, IconPhotoDown, IconRestore, IconUpload, IconX } from '@tabler/icons-react'
import { Dropzone, FileWithPath, IMAGE_MIME_TYPE } from '@mantine/dropzone'

interface ImageDropzoneWithPreviewProps {
  image: FileWithPath | string | null
  setImage: Dispatch<SetStateAction<FileWithPath | string | null>>
  defaultValue?: string | null
  label?: string
  ariaLabel?: string
}

function LargeImageDropzoneWithPreview({
  image,
  setImage,
  defaultValue,
  label = 'Image',
  ariaLabel = 'image'
}: ImageDropzoneWithPreviewProps) {
  function handleRemoveImage() {
    setImage(null)
  }

  function handleResetImage() {
    setImage(defaultValue)
  }

  if (image) {
    return (
      <Group justify={'center'}>
        <AspectRatio w={'50%'} ratio={1}>
          <Image
            src={typeof image === 'string' ? image : URL.createObjectURL(image)}
            radius={'md'}
            alt={`${ariaLabel}-preview`}
          />
        </AspectRatio>

        <Stack gap={'xs'}>
          {defaultValue && (
            <Tooltip
              label={
                image === defaultValue
                  ? `Original ${label} can be restored after re-upload`
                  : `Reset ${label}`
              }
              openDelay={300}
              position={'right'}
            >
              <ActionIcon
                variant={'subtle'}
                aria-label={`reset-${ariaLabel}`}
                size={'lg'}
                disabled={image === defaultValue}
                onClick={handleResetImage}
              >
                <IconRestore size={18} />
              </ActionIcon>
            </Tooltip>
          )}

          <FileButton
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-ignore
            inputProps={{ 'data-testid': `upload-${ariaLabel}-input` }}
            onChange={setImage}
            accept={IMAGE_MIME_TYPE.join(',')}
          >
            {(props) => (
              <Tooltip label={`Upload another ${label}`} openDelay={300} position={'right'}>
                <ActionIcon
                  c={'dark'}
                  variant={'subtle'}
                  aria-label={`upload-${ariaLabel}`}
                  size={'lg'}
                  {...props}
                >
                  <IconPhotoDown size={18} />
                </ActionIcon>
              </Tooltip>
            )}
          </FileButton>

          <Tooltip label={`Remove ${label}`} openDelay={300} position={'right'}>
            <ActionIcon
              variant={'subtle'}
              color={'red.6'}
              aria-label={`remove-${ariaLabel}`}
              size={'lg'}
              sx={(theme) => ({
                '&:hover': { backgroundColor: alpha(theme.colors.red[2], 0.7) }
              })}
              onClick={handleRemoveImage}
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
      aria-label={`${ariaLabel}-dropzone`}
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      inputProps={{ 'data-testid': `${ariaLabel}-dropzone-input` }}
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

        <Text fz={'xl'} c={'inherit'}>
          Drag {label} here or click to select it
        </Text>
      </Group>
    </Dropzone>
  )
}

export default LargeImageDropzoneWithPreview
