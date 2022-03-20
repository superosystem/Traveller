package com.gusrylmubarok.springbasic.service;

import com.gusrylmubarok.springbasic.repository.CategoryRepository;
import lombok.Getter;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Component
public class CategoryService {

    @Getter
    private CategoryRepository categoryRepository;

    @Autowired
    public void setCategoryRepository(CategoryRepository categoryRepository) {
        this.categoryRepository = categoryRepository;
    }
}
