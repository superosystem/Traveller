package com.gusrylmubarok.spring.perpustakaan.entity;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.Table;

@Entity
@Table(name = "book")
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
public class Book extends BaseEntity<Book> {

    private static final long serialVersionUID = -397387408185673271L;

    @Column(name = "title", columnDefinition = "VARCHAR(100)")
    private String title;

    @Column(name = "isbn")
    private String isbn;

    @Column(name = "author")
    private String author;

    @Column(name = "publisher")
    private String publisher;

    @Column(name = "description")
    private String description;

}
