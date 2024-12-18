import { Dispatch, ReactElement, SetStateAction } from 'react'
import { ActionIcon, alpha, Box, FileButton, Group, Image, Tooltip } from '@mantine/core'
import { IconMusic, IconPhotoDown, IconUpload, IconX } from '@tabler/icons-react'
import { Dropzone, FileWithPath, IMAGE_MIME_TYPE } from '@mantine/dropzone'

interface ImageDropzoneWithPreviewProps {
  image: FileWithPath
  setImage: Dispatch<SetStateAction<FileWithPath>>
  w?: number
  h?: number
  icon?: ReactElement
  iconSizes?: number
}

function ImageDropzoneWithPreview({
  image,
  setImage,
  icon,
  w = 92,
  h = 92,
  iconSizes = 45
}: ImageDropzoneWithPreviewProps) {
  if (image) {
    return (
      <Box pos={'relative'}>
        <Image src={URL.createObjectURL(image)} w={w} h={h} radius={'32px'} alt={'image-preview'} />

        <Box pos={'absolute'} top={h - 22} left={-8}>
          <Tooltip label={'Remove Image'} openDelay={300} position={'bottom'}>
            <ActionIcon
              c={'white'}
              radius={'50%'}
              aria-label={'remove-image-button'}
              size={'md'}
              sx={(theme) => ({
                transition: '0.15s',
                backgroundColor: alpha(theme.colors.red[5], 0.5),
                '&:hover': { backgroundColor: alpha(theme.colors.red[5], 0.7) }
              })}
              onClick={() => setImage(null)}
            >
              <IconX size={15} />
            </ActionIcon>
          </Tooltip>
        </Box>
        <Box pos={'absolute'} top={h - 22} right={-8}>
          <FileButton onChange={setImage} accept="image/png,image/jpeg">
            {(props) => (
              <Tooltip label={'Reload Image'} openDelay={300} position={'bottom'}>
                <ActionIcon
                  c={'dark'}
                  radius={'50%'}
                  aria-label={'add-image-button'}
                  size={'md'}
                  {...props}
                  sx={(theme) => ({
                    transition: '0.15s',
                    backgroundColor: alpha(theme.colors.gray[4], 0.5),
                    '&:hover': {
                      backgroundColor: alpha(theme.colors.gray[4], 0.7)
                    }
                  })}
                >
                  <IconPhotoDown size={15} />
                </ActionIcon>
              </Tooltip>
            )}
          </FileButton>
        </Box>
      </Box>
    )
  }

  return (
    <Dropzone
      onDrop={(files) => setImage(files[0])}
      accept={IMAGE_MIME_TYPE}
      multiple={false}
      w={w}
      h={h}
      styles={{
        inner: {
          height: 'calc(100% - 5px)'
        }
      }}
      sx={(theme) => ({
        cursor: 'pointer',
        borderRadius: '32px',
        border: `1px solid ${theme.colors.gray[4]}`,
        transition: '0.3s',
        color: theme.colors.gray[6],
        backgroundColor: theme.colors.gray[3],
        '&:hover': {
          color: theme.colors.gray[7],
          backgroundColor: theme.colors.gray[4]
        },

        '&:where([data-accept])': {
          color: theme.colors.green[7],
          backgroundColor: theme.colors.green[3]
        },

        '&:where([data-reject])': {
          color: theme.colors.red[7],
          backgroundColor: theme.colors.red[3]
        }
      })}
    >
      <Group justify="center" style={{ pointerEvents: 'none' }} h={'100%'}>
        <Dropzone.Accept>
          <IconUpload size={iconSizes} />
        </Dropzone.Accept>
        <Dropzone.Reject>
          <IconX size={iconSizes} />
        </Dropzone.Reject>
        <Dropzone.Idle>{icon ? icon : <IconMusic size={iconSizes} />}</Dropzone.Idle>
      </Group>
    </Dropzone>
  )
}

export default ImageDropzoneWithPreview