# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# Note that this schema.rb definition is the authoritative source for your
# database schema. If you need to create the application database on another
# system, you should be using db:schema:load, not running all the migrations
# from scratch. The latter is a flawed and unsustainable approach (the more migrations
# you'll amass, the slower it'll run and the greater likelihood for issues).
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema.define(version: 20210214191413) do

  create_table "apps", force: :cascade, options: "ENGINE=InnoDB DEFAULT CHARSET=latin1" do |t|
    t.string   "token"
    t.string   "name"
    t.integer  "chats_count", default: 0
    t.datetime "created_at",              null: false
    t.datetime "updated_at",              null: false
    t.index ["token"], name: "index_apps_on_token", using: :btree
  end

  create_table "chats", force: :cascade, options: "ENGINE=InnoDB DEFAULT CHARSET=latin1" do |t|
    t.integer  "number"
    t.string   "name"
    t.integer  "messages_count", default: 0
    t.integer  "app_id"
    t.datetime "created_at",                 null: false
    t.datetime "updated_at",                 null: false
    t.index ["app_id"], name: "index_chats_on_app_id", using: :btree
    t.index ["number"], name: "index_chats_on_number", using: :btree
  end

  create_table "messages", force: :cascade, options: "ENGINE=InnoDB DEFAULT CHARSET=latin1" do |t|
    t.integer  "number"
    t.text     "content",    limit: 65535
    t.integer  "chat_id"
    t.datetime "created_at",               null: false
    t.datetime "updated_at",               null: false
    t.index ["chat_id"], name: "index_messages_on_chat_id", using: :btree
    t.index ["number"], name: "index_messages_on_number", using: :btree
  end

end
