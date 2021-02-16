Rails.application.routes.draw do
  # For details on the DSL available within this file, see http://guides.rubyonrails.org/routing.html
  scope :api do
    resources :apps, path: :applications, only: [:index, :show, :update], param: :token do 
      resources :chats, only: [:index, :show, :update], param: :number do 
        get 'search', to: 'chats#search'
        resources :messages, only: [:index, :show, :update], param: :number
      end
    end
  end
end
