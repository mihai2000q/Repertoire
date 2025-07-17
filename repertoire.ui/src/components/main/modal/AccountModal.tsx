import {
  Button,
  Group,
  LoadingOverlay,
  Modal,
  ScrollArea,
  Stack,
  Text,
  TextInput,
  Tooltip
} from '@mantine/core'
import LargeImageDropzoneWithPreview from '../../@ui/image/LargeImageDropzoneWithPreview.tsx'
import { useEffect, useState } from 'react'
import { useForm } from '@mantine/form'
import { zod4Resolver } from 'mantine-form-zod-resolver'
import User from '../../../types/models/User.ts'
import { AccountForm, accountSchema } from '../../../validation/mainForm.ts'
import { toast } from 'react-toastify'
import { useDidUpdate } from '@mantine/hooks'
import { FileWithPath } from '@mantine/dropzone'
import dayjs from 'dayjs'
import {
  useDeleteProfilePictureMutation,
  useSaveProfilePictureMutation,
  useUpdateUserMutation
} from '../../../state/api/usersApi.ts'

interface AccountModalProps {
  opened: boolean
  onClose: () => void
  user: User
}

function AccountModal({ opened, onClose, user }: AccountModalProps) {
  const [updateUserMutation, { isLoading: isUpdateLoading }] = useUpdateUserMutation()
  const [saveProfilePictureMutation, { isLoading: isSaveProfilePictureLoading }] =
    useSaveProfilePictureMutation()
  const [deleteProfilePictureMutation, { isLoading: isDeleteProfilePictureLoading }] =
    useDeleteProfilePictureMutation()
  const isLoading = isUpdateLoading || isSaveProfilePictureLoading || isDeleteProfilePictureLoading

  const [userHasChanged, setUserHasChanged] = useState(false)
  const [pictureHasChanged, setPictureHasChanged] = useState(false)
  const hasChanged = userHasChanged || pictureHasChanged

  const form = useForm<AccountForm>({
    mode: 'uncontrolled',
    initialValues: {
      name: user.name,
      profilePicture: user.profilePictureUrl
    },
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zod4Resolver(accountSchema),
    onValuesChange: (values) => {
      setUserHasChanged(values.name !== user.name)
      setPictureHasChanged(values.profilePicture !== user.profilePictureUrl)
    }
  })

  const [profilePicture, setProfilePicture] = useState<string | FileWithPath>(
    user.profilePictureUrl
  )
  useEffect(() => form.setFieldValue('profilePicture', profilePicture), [profilePicture])
  useDidUpdate(() => setProfilePicture(user.profilePictureUrl), [user])

  async function updateUser({ name, profilePicture }: AccountForm) {
    if (userHasChanged)
      await updateUserMutation({
        name: name.trim()
      }).unwrap()

    if (profilePicture !== null && typeof profilePicture !== 'string')
      await saveProfilePictureMutation({
        profile_pic: profilePicture as FileWithPath
      })
    else if (profilePicture === null && user.profilePictureUrl) await deleteProfilePictureMutation()

    toast.info('Account updated!')
    onClose()
    setUserHasChanged(false)
    setPictureHasChanged(false)
  }

  return (
    <Modal opened={opened} onClose={onClose} title={'Account'} styles={{ body: { padding: 0 } }}>
      <ScrollArea.Autosize offsetScrollbars={'y'} scrollbars={'y'} scrollbarSize={7} mah={'77vh'}>
        <LoadingOverlay visible={isLoading} loaderProps={{ type: 'bars' }} />

        <form onSubmit={form.onSubmit(updateUser)}>
          <Stack px={26} pb={'md'}>
            <LargeImageDropzoneWithPreview
              image={profilePicture}
              setImage={setProfilePicture}
              defaultValue={user.profilePictureUrl}
              label={'Picture'}
              ariaLabel={'profile-picture'}
            />

            <TextInput
              withAsterisk={true}
              maxLength={100}
              label="Name"
              placeholder="Your name"
              key={form.key('name')}
              {...form.getInputProps('name')}
            />

            <TextInput label="Email" disabled={true} defaultValue={user.email} />

            <Group justify={'space-between'}>
              <Text fz={'sm'} fw={500} c={'dimmed'} inline>
                Created on <b>{dayjs(user.createdAt).format('DD MMM YYYY')}</b>
              </Text>

              {user.createdAt !== user.updatedAt && (
                <Text fz={'sm'} fw={500} c={'dimmed'} inline>
                  Last Modified on <b>{dayjs(user.updatedAt).format('DD MMM YYYY')}</b>
                </Text>
              )}
            </Group>

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
      </ScrollArea.Autosize>
    </Modal>
  )
}

export default AccountModal
