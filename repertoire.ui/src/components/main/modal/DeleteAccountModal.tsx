import {
  Box,
  Button,
  Group,
  Loader,
  LoadingOverlay,
  Modal,
  PasswordInput,
  Stack,
  Text,
  TextInput,
  Transition
} from '@mantine/core'
import { useState } from 'react'
import { signOut } from '../../../state/slice/authSlice.ts'
import { useAppDispatch } from '../../../state/store.ts'
import { useNavigate } from 'react-router-dom'
import HttpErrorResponse from '../../../types/responses/HttpErrorResponse.ts'
import User from '../../../types/models/User.ts'
import { useForm } from '@mantine/form'
import { zod4Resolver } from 'mantine-form-zod-resolver'
import { DeleteAccountForm, deleteAccountSchema } from '../../../validation/mainForm.ts'
import { useSignInMutation } from '../../../state/authApi.ts'
import { useDeleteUserMutation } from '../../../state/api/usersApi.ts'
import { toast } from 'react-toastify'

interface DeleteAccountModalProps {
  opened: boolean
  onClose: () => void
  onCloseSettingsModal: () => void
  user: User
}

function DeleteAccountModal({
  opened,
  onClose,
  onCloseSettingsModal,
  user
}: DeleteAccountModalProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const [deleteMutation, { isLoading: isDeleteLoading }] = useDeleteUserMutation()
  const [signInMutation, { error, isLoading: isSignInLoading }] = useSignInMutation()
  const signInError = (error as HttpErrorResponse | undefined)?.data?.error
  const isLoading = isDeleteLoading || isSignInLoading

  const [activeStep, setActiveStep] = useState(1)

  const form = useForm<DeleteAccountForm>({
    mode: 'uncontrolled',
    initialValues: {
      password: ''
    },
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zod4Resolver(deleteAccountSchema)
  })

  const onCloseWithStep = () => {
    onClose()
    setActiveStep(1)
  }

  function handleContinue() {
    setActiveStep(2)
  }

  async function handleDelete({ password }: DeleteAccountForm) {
    try {
      await signInMutation({ email: user.email, password }).unwrap()
    } catch (e) {
      return
    }
    await deleteMutation()
    onCloseWithStep()
    onCloseSettingsModal()
    dispatch(signOut())
    navigate('sign-in')
    toast(`We are sorry to let you down. Goodbye! :(`)
  }

  return (
    <Modal
      opened={opened}
      onClose={onCloseWithStep}
      title={'Delete Account'}
      overlayProps={{ blur: 1 }}
      closeOnClickOutside={false}
      centered
    >
      <LoadingOverlay
        visible={isLoading}
        loaderProps={{
          children: (
            <Stack align={'center'}>
              <Loader />
              {isSignInLoading && (
                <Text fw={500} c={'dimmed'}>
                  Authenticating...
                </Text>
              )}
              {isDeleteLoading && (
                <Text fw={500} c={'dimmed'}>
                  Deleting all resources...
                </Text>
              )}
            </Stack>
          )
        }}
      />

      <Box pos={'relative'} py={0}>
        <Transition
          mounted={activeStep === 1}
          transition="fade-left"
          duration={400}
          timingFunction="ease"
        >
          {(styles) => (
            <Stack style={styles} pos={activeStep === 1 ? 'unset' : 'absolute'}>
              <Stack gap={'xxs'}>
                <Text fw={500}>Are you sure you want to delete your account?</Text>
                <Text fz={'xs'} fw={500} c={'dimmed'}>
                  This action will result in the immediate loss of access to Repertoire and the
                  permanent removal of your account data. There will be no option for recovery.
                </Text>
              </Stack>

              <Group gap={'xxs'} style={{ alignSelf: 'end' }}>
                <Button variant={'subtle'} onClick={onCloseWithStep}>
                  Keep Account
                </Button>
                <Button onClick={handleContinue}>Continue</Button>
              </Group>
            </Stack>
          )}
        </Transition>

        <Transition
          mounted={activeStep === 2}
          transition="fade-right"
          duration={400}
          timingFunction="ease"
        >
          {(styles) => (
            <form onSubmit={form.onSubmit(handleDelete)}>
              <Stack style={styles} pos={activeStep === 2 ? 'unset' : 'absolute'}>
                <TextInput
                  label="Email"
                  placeholder="Your email"
                  value={user.email}
                  disabled={true}
                />
                <PasswordInput
                  role={'textbox'}
                  label="Password"
                  placeholder="Your password"
                  key={form.key('password')}
                  {...form.getInputProps('password')}
                  {...(signInError && { error: signInError })}
                  disabled={isLoading}
                />

                <Group gap={'xxs'} style={{ alignSelf: 'end' }}>
                  <Button variant={'subtle'} onClick={onCloseWithStep}>
                    Cancel
                  </Button>
                  <Button bg={'red'} type={'submit'}>
                    Delete
                  </Button>
                </Group>
              </Stack>
            </form>
          )}
        </Transition>
      </Box>
    </Modal>
  )
}

export default DeleteAccountModal
