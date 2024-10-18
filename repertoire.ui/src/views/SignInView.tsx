import {
  Anchor,
  Button,
  Container,
  Flex,
  Paper,
  PasswordInput,
  Stack,
  Text,
  TextInput,
  Title
} from '@mantine/core'
import { ReactElement } from 'react'
import { useSignInMutation } from '../state/api'
import { useAppDispatch } from '../state/store'
import { setToken } from '../state/authSlice'
import HttpErrorResponse from '../types/responses/HttpErrorResponse'
import { useLocation, useNavigate } from 'react-router-dom'
import { useForm, zodResolver } from '@mantine/form'
import { SignInForm, signInValidation } from '../validation/signInForm'

function SignInView(): ReactElement {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const location = useLocation()

  const [signInMutation, { error, isLoading }] = useSignInMutation()
  const loginError = (error as HttpErrorResponse | undefined)?.data?.error

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      email: '',
      password: ''
    } as SignInForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(signInValidation)
  })

  async function signIn({ email, password }: SignInForm): Promise<void> {
    try {
      const res = await signInMutation({ email, password }).unwrap()

      dispatch(setToken(res.token))
      navigate(location.state?.from?.pathname ?? 'home')
    } catch (e) {
      /*ignored*/
    }
  }

  return (
    <Container my={40}>
      <Flex direction={'column'} align={'center'} justify={'center'} h={'90%'}>
        <Title ta="center" order={2}>
          Welcome back!
        </Title>
        <Text c="dimmed" size="sm" ta="center" mt={5}>
          Do not have an account yet?{' '}
          <Anchor size="sm" component="button">
            Create account
          </Anchor>
        </Text>

        <Paper withBorder shadow="md" p={30} mt={15}>
          <form onSubmit={form.onSubmit(signIn)}>
            <Stack align={'flex-start'} gap={0} w={200}>
              <Stack w={'100%'}>
                <TextInput
                  required
                  label="Email"
                  placeholder="Your email"
                  key={form.key('email')}
                  {...form.getInputProps('email')}
                  {...(loginError && { error: loginError })}
                  maxLength={256}
                  disabled={isLoading}
                />
                <PasswordInput
                  required
                  label="Password"
                  placeholder="Your password"
                  key={form.key('password')}
                  {...form.getInputProps('password')}
                  {...(loginError && { error: loginError })}
                  disabled={isLoading}
                />
              </Stack>
              <Anchor component="button" size="sm" mt={6} style={{ alignSelf: 'flex-end' }}>
                Forgot password?
              </Anchor>
              <Button type={'submit'} fullWidth mt={'sm'} disabled={isLoading}>
                Sign in
              </Button>
            </Stack>
          </form>
        </Paper>
      </Flex>
    </Container>
  )
}

export default SignInView