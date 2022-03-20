package com.gusrylmubarok.springbasic.service;

import com.gusrylmubarok.springbasic.repository.ProductRepository;
import lombok.Getter;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

// @Scope("prototype")
// @Lazy
@Component
public class ProductService {

    @Getter
    private ProductRepository productRepository;

    @Autowired
    public ProductService(ProductRepository productRepository){
        this.productRepository = productRepository;
    }

    public ProductService(ProductRepository productRepository, String name){
        this.productRepository = productRepository;
    }

}
