class AppsController < ApplicationController
  before_action :set_app, only: [:show, :update]

  def index
    @apps = App.all
    render json: @apps
  end

  def show
    render json: @app
  end

  def update
    if @app.update(name: app_params['application_name'])
      render json: @app
    else
      render json: @app.errors, status: :unprocessable_entity
    end
  end

  private
    # Use callbacks to share common setup or constraints between actions.
    def set_app
      @app = App.find_by(token: params[:token])
    end

    # Only allow a trusted parameter "white list" through.
    def app_params
      params.permit('application_name')
    end
end
