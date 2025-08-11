import { mantineRender } from '../../../test-utils.tsx'
import HomeRecentlyPlayedCard from './HomeRecentlyPlayedCard.tsx'
import { expect } from 'vitest'
import { screen } from '@testing-library/react'
import dayjs from 'dayjs'
import { userEvent } from '@testing-library/user-event'

describe('Home Recently Played Card', () => {
  it('should render', async () => {
    const user = userEvent.setup()

    const imageUrl = 'something.png'
    const title = 'Song 1'
    const progress = 12
    const lastPlayed = '2022-11-11T10:30'
    const openedMenu = false
    const onClick = vi.fn()

    mantineRender(
      <HomeRecentlyPlayedCard
        imageUrl={imageUrl}
        title={title}
        progress={progress}
        lastPlayed={lastPlayed}
        openedMenu={openedMenu}
        defaultIcon={<></>}
        onClick={onClick}
      />
    )

    expect(screen.getByRole('img', { name: title })).toHaveAttribute('src', imageUrl)
    expect(screen.getByText(title)).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'progress' })).toBeInTheDocument()
    expect(screen.getByText(dayjs(lastPlayed).format('DD MMM'))).toBeInTheDocument()

    await user.hover(screen.getByText(dayjs(lastPlayed).format('DD MMM')))
    expect(
      await screen.findByText(new RegExp(dayjs(lastPlayed).format('D MMMM YYYY [at] hh:mm A'), 'i'))
    ).toBeInTheDocument()

    await user.click(screen.getByText(title))
    expect(onClick).toHaveBeenCalledOnce()
  })

  it('should render with default icon instead of image', () => {
    const defaultIconTestId = 'default-icon-testid'
    const defaultIcon = <div data-testid={defaultIconTestId}>Def</div>

    mantineRender(
      <HomeRecentlyPlayedCard
        imageUrl={null}
        title={'Song 1'}
        progress={12}
        lastPlayed={'2022-11-11T10:30'}
        openedMenu={false}
        defaultIcon={defaultIcon}
        onClick={vi.fn()}
      />
    )

    expect(screen.getByTestId(defaultIconTestId)).toBeInTheDocument()
  })

  it('should render with default icon instead of image', async () => {
    const user = userEvent.setup()

    const additionalTextContent = 'some-text'
    const additionalTextOnClick = vi.fn()

    mantineRender(
      <HomeRecentlyPlayedCard
        imageUrl={null}
        title={'Song 1'}
        progress={12}
        lastPlayed={'2022-11-11T10:30'}
        openedMenu={false}
        defaultIcon={<></>}
        onClick={vi.fn()}
        additionalText={{
          content: additionalTextContent,
          onClick: additionalTextOnClick
        }}
      />
    )

    expect(screen.getByText(additionalTextContent)).toBeInTheDocument()
    expect(screen.getByText(additionalTextContent)).toBeInTheDocument()

    await user.click(screen.getByText(additionalTextContent))
    expect(additionalTextOnClick).toHaveBeenCalledOnce()
  })
})
