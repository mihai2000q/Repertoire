import { Button, LoadingOverlay, Modal, Stack, TextInput, Tooltip } from '@mantine/core'
import LargeImageDropzoneWithPreview from '../../@ui/image/LargeImageDropzoneWithPreview.tsx'
import { useEffect, useState } from 'react'
import { useForm, zodResolver } from '@mantine/form'
import User from '../../../types/models/User.ts'
import { AccountForm, accountValidation } from '../../../validation/mainForm.ts'
import {
  useDeleteProfilePictureFromUserMutation,
  useSaveProfilePictureToUserMutation,
  useUpdateUserMutation
} from '../../../state/api.ts'
import { toast } from 'react-toastify'
import { useDidUpdate } from '@mantine/hooks'
import { FileWithPath } from '@mantine/dropzone'

interface AccountModalProps {
  opened: boolean
  onClose: () => void
  user: User
}

function AccountModal({ opened, onClose, user }: AccountModalProps) {
  const [updateUserMutation, { isLoading: isUpdateLoading }] = useUpdateUserMutation()
  const [saveProfilePictureMutation, { isLoading: isSaveProfilePictureLoading }] =
    useSaveProfilePictureToUserMutation()
  const [deleteProfilePictureMutation, { isLoading: isDeleteProfilePictureLoading }] =
    useDeleteProfilePictureFromUserMutation()
  const isLoading = isUpdateLoading || isSaveProfilePictureLoading || isDeleteProfilePictureLoading

  const [hasChanged, setHasChanged] = useState(false)

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      name: user.name,
      profilePicture: user.profilePictureUrl
    } as AccountForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(accountValidation),
    onValuesChange: (values) => {
      setHasChanged(values.name !== user.name || values.profilePicture !== user.profilePictureUrl)
    }
  })

  const [profilePicture, setProfilePicture] = useState<string | FileWithPath>(
    user.profilePictureUrl
  )
  useEffect(() => form.setFieldValue('profilePicture', profilePicture), [profilePicture])
  useDidUpdate(() => setProfilePicture(user.profilePictureUrl), [user])

  async function updateUser({ name, profilePicture }: AccountForm) {
    name = name.trim()

    await updateUserMutation({
      name: name
    }).unwrap()

    if (profilePicture !== null && typeof profilePicture !== 'string') {
      await saveProfilePictureMutation({
        profile_pic: profilePicture
      })
    } else if (profilePicture === null && user.profilePictureUrl) {
      await deleteProfilePictureMutation()
    }

    toast.info('Account updated!')
    onClose()
    setHasChanged(false)
  }

  return (
    <Modal opened={opened} onClose={onClose} title={'Account'}>
      <Modal.Body px={'xs'} py={0}>
        <LoadingOverlay visible={isLoading} />

        <form onSubmit={form.onSubmit(updateUser)}>
          <Stack>
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
      </Modal.Body>
    </Modal>
  )
}

export default AccountModal
