import { Dispatch, ReactElement, SetStateAction, useState } from 'react'
import {
  ActionIcon,
  alpha,
  AspectRatio,
  Box,
  FileButton,
  Group,
  Image,
  Menu,
  Stack,
  Tooltip
} from '@mantine/core'
import { IconMusic, IconPhotoFilled, IconTrashFilled, IconUpload, IconX } from '@tabler/icons-react'
import { Dropzone, FileWithPath, IMAGE_MIME_TYPE } from '@mantine/dropzone'

interface ImageDropzoneWithPreviewProps {
  image: FileWithPath | null
  setImage: Dispatch<SetStateAction<FileWithPath | null>>
  w?: number
  h?: number
  icon?: ReactElement
  iconSizes?: number
  radius?: string
}

function ImageDropzoneWithPreview({
  image,
  setImage,
  icon,
  w = 92,
  h = 92,
  iconSizes = 45,
  radius = '24px'
}: ImageDropzoneWithPreviewProps) {
  const [isMenuOpened, setIsMenuOpened] = useState(false)

  function handleRemoveImage() {
    setImage(null)
  }

  function handleImageChange(image: FileWithPath) {
    setImage(image)
    setIsMenuOpened(false)
  }

  if (image) {
    return (
      <Box pos={'relative'}>
        <AspectRatio>
          <Image
            src={URL.createObjectURL(image)}
            w={w}
            h={h}
            radius={radius}
            alt={'image-preview'}
          />
        </AspectRatio>

        <Box pos={'absolute'} top={0} right={0} h={h} w={h}>
          <Stack gap={0} align={'center'} h={'100%'} justify={'center'}>
            <Menu
              opened={isMenuOpened}
              onChange={setIsMenuOpened}
              withArrow
              offset={-iconSizes / 2}
              transitionProps={{ transition: 'fade-up', duration: 150 }}
            >
              <Menu.Target>
                <Tooltip label={'Open image options menu'} openDelay={500}>
                  <ActionIcon
                    aria-label={'image-options'}
                    h={'100%'}
                    w={'100%'}
                    radius={radius}
                    style={(theme) => ({
                      transition: '0.3s',
                      backgroundColor: alpha(theme.black, 0.5),
                      color: theme.white
                    })}
                    sx={{
                      opacity: isMenuOpened ? 1 : 0,
                      '&:hover': { opacity: 1 }
                    }}
                  >
                    <IconPhotoFilled size={iconSizes / 1.2} />
                  </ActionIcon>
                </Tooltip>
              </Menu.Target>

              <Menu.Dropdown>
                <FileButton
                  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                  // @ts-ignore
                  inputProps={{ 'data-testid': 'upload-image-input' }}
                  onChange={handleImageChange}
                  accept={IMAGE_MIME_TYPE.join(',')}
                >
                  {(props) => (
                    <Menu.Item
                      leftSection={<IconUpload size={18} />}
                      closeMenuOnClick={false}
                      {...props}
                    >
                      Upload Image
                    </Menu.Item>
                  )}
                </FileButton>
                <Menu.Item
                  c={'red'}
                  leftSection={<IconTrashFilled size={18} />}
                  onClick={handleRemoveImage}
                >
                  Remove Image
                </Menu.Item>
              </Menu.Dropdown>
            </Menu>
          </Stack>
        </Box>
      </Box>
    )
  }

  return (
    <Dropzone
      aria-label={'image-dropzone'}
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      inputProps={{ 'data-testid': 'image-dropzone-input' }}
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
        borderRadius: radius,
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
