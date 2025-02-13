import {mantineRender} from "../../test-utils.tsx";
import HomeGenres from "./HomeGenres.tsx";
import {screen} from "@testing-library/react";

describe('Home Genres', () => {
  it('should render', () => {
    mantineRender(<HomeGenres />)

    expect(screen.getByText(/genres/i)).toBeInTheDocument()
    expect(screen.getByText(/see all/i)).toBeInTheDocument()
  })
})
