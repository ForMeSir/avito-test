package repository

import (
	"avito/internal/banner"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type BannerPostgres struct {
	db *sqlx.DB
}
const(
	table="banners"
)
func NewBannerPostgres(db *sqlx.DB) *BannerPostgres{
	return &BannerPostgres{db:db}
}

func(r *BannerPostgres) FindOne(tag int, feature int) (banner banner.Content, err error){
	query := fmt.Sprintf("SELECT content_title, content_text, content_url, is_active FROM %s WHERE feature_id=$1 AND $2 = ANY(tag_ids)", table)
	err =r.db.Get(&banner,query, feature,tag)
	return
}

func(r *BannerPostgres) Create(banner banner.Banner)(id int,err error){
	query := fmt.Sprintf("INSERT INTO %s (tag_ids,feature_id,content_title,content_text,content_url,is_active,created_at,updated_at) values($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id", table)
	row:=r.db.QueryRow(query,pq.Array(banner.TagIds),banner.FeatureId,banner.Content.Title,banner.Content.Text,banner.Content.Url,banner.Content.IsActive,time.Now().Format(time.RFC3339),time.Now().Format(time.RFC3339))
	if err=row.Scan(&id); err!=nil{
		fmt.Println(err)
		return
	}
	return
}

func(r *BannerPostgres) FindAllByFilter(feature int, tag int, limit int, offset int) (ban []banner.FullBanner, err error){
	var rows *sql.Rows
	if tag==0 && feature==0 {
		  query:=fmt.Sprintf("SELECT * FROM %s LIMIT $1 OFFSET $2",table)

				if limit==0{
					query=fmt.Sprintf("SELECT * FROM %s LIMIT ALL OFFSET $1",table)
					rows,err=r.db.Query(query,offset)
				}else {
					rows,err=r.db.Query(query,limit,offset)
				}

	}else {
		query:=fmt.Sprintf("SELECT * FROM %s WHERE feature_id=$1 OR $2 = ANY(tag_ids) LIMIT $3 OFFSET $4", table)
		  
		  if limit==0{
				query=fmt.Sprintf("SELECT * FROM %s WHERE feature_id=$1 OR $2 = ANY(tag_ids) LIMIT ALL OFFSET $3",table)
				rows,err=r.db.Query(query,feature,tag,offset)
			}else {
				rows,err=r.db.Query(query,feature,tag,limit,offset)
			}

		}
				for rows.Next() {
					var banOne banner.FullBanner
					if err := rows.Scan(&banOne.Id,&banOne.TagIds, &banOne.FeatureId, &banOne.Title,&banOne.Text,&banOne.Url,&banOne.IsActive,&banOne.CreatedAt, &banOne.UpdatedAt); err != nil {
							log.Fatal(err)
					}
					ban = append(ban, banOne)
		   	}
				return

}

func(r *BannerPostgres) Update(id int,update banner.UpdateBody)(err error){
	if update.TagIds.Valid{
		query:=fmt.Sprintf("UPDATE %s SET tag_ids = $1 WHERE id = $2", table)
	  _,err=r.db.Exec(query,pq.Array(update.TagIds.Val),id)
	  if err!=nil{
			return
		}
	}

  if update.FeatureId.Valid{
		query:=fmt.Sprintf("UPDATE %s SET feature_id = $1 WHERE id = $2", table)
	  _,err=r.db.Exec(query,update.FeatureId.Val,id)
	  if err!=nil{
			return
		}
	}

	if update.Content.Valid{
		query:=fmt.Sprintf("UPDATE %s SET content_title = $1, content_text=$2, content_url=$3 WHERE id = $4", table)
	  _,err=r.db.Exec(query,update.Content.Val.Title, update.Content.Val.Text, update.Content.Val.Url, id )
	  if err!=nil{
			return
		}
	}

	if update.IsActive.Valid{
		query:=fmt.Sprintf("UPDATE %s SET is_active = $1 WHERE id = $2", table)
	  _,err=r.db.Exec(query,update.IsActive.Val,id)
	  if err!=nil{
			return
		}
	}

	query:=fmt.Sprintf("UPDATE %s SET updated_at=$1 WHERE id = $2", table)
	_,err=r.db.Exec(query,time.Now().Format(time.RFC3339),id)
	if err!=nil{
		return
	}
	return nil
}


func(r *BannerPostgres) Delete(id int) error{
	query:=fmt.Sprintf("DELETE FROM %s WHERE id = $1", table)
	  _,err:=r.db.Exec(query,id )
	  if err!=nil{
			return err
		}
	return nil
}