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
import { useLocation, useNavigate } from 'react-router-dom'
import { useForm, zodResolver } from '@mantine/form'
import { SignInForm, signInValidation } from '../validation/signInForm'
import {useSignInMutation} from "../state/signInApi";

function SignInView(): ReactElement {
  const navigate = useNavigate()
  const location = useLocation()

  const { data: signInResponse, mutate, error, isPending } = useSignInMutation()
  const signInError = error.message

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
    mutate({ email, password })
    localStorage.setItem('token', signInResponse.data.token)
    navigate(location.state?.from?.pathname ?? 'home')
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
                  label="Email"
                  placeholder="Your email"
                  key={form.key('email')}
                  {...form.getInputProps('email')}
                  {...(signInError && { error: signInError })}
                  maxLength={256}
                  disabled={isPending}
                />
                <PasswordInput
                  label="Password"
                  placeholder="Your password"
                  key={form.key('password')}
                  {...form.getInputProps('password')}
                  {...(signInError && { error: signInError })}
                  disabled={isPending}
                />
              </Stack>
              <Anchor component="button" size="sm" mt={6} style={{ alignSelf: 'flex-end' }}>
                Forgot password?
              </Anchor>
              <Button type={'submit'} fullWidth mt={'sm'} disabled={isPending}>
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
